package api

import (
	"bytes"
	"context"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// ObjectStore backs Surfing media bytes. Prefer direct-to-R2 drafts; PVC is legacy fallback.
type ObjectStore interface {
	Enabled() bool
	Put(ctx context.Context, key string, body []byte, contentType string) error
	Delete(ctx context.Context, key string) error
	PublicURL(key string) string
	PresignPut(ctx context.Context, key, contentType string, expiry time.Duration) (string, error)
	Copy(ctx context.Context, srcKey, dstKey, contentType string) error
	Head(ctx context.Context, key string) (size int64, ok bool, err error)
}

type noopStore struct{}

func (noopStore) Enabled() bool { return false }
func (noopStore) Put(context.Context, string, []byte, string) error {
	return fmt.Errorf("object store disabled")
}
func (noopStore) Delete(context.Context, string) error { return nil }
func (noopStore) PublicURL(string) string              { return "" }
func (noopStore) PresignPut(context.Context, string, string, time.Duration) (string, error) {
	return "", fmt.Errorf("object store disabled")
}
func (noopStore) Copy(context.Context, string, string, string) error {
	return fmt.Errorf("object store disabled")
}
func (noopStore) Head(context.Context, string) (int64, bool, error) { return 0, false, nil }

type r2Store struct {
	client     *s3.Client
	presigner  *s3.PresignClient
	bucket     string
	publicBase string
}

func (r *r2Store) Enabled() bool { return r != nil && r.client != nil && r.bucket != "" }

func (r *r2Store) Put(ctx context.Context, key string, body []byte, contentType string) error {
	_, err := r.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:       aws.String(r.bucket),
		Key:          aws.String(key),
		Body:         bytes.NewReader(body),
		ContentType:  aws.String(contentType),
		CacheControl: aws.String("public, max-age=31536000, immutable"),
	})
	return err
}

func (r *r2Store) Delete(ctx context.Context, key string) error {
	_, err := r.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(key),
	})
	return err
}

func (r *r2Store) PublicURL(key string) string {
	base := strings.TrimRight(r.publicBase, "/")
	if base == "" {
		return ""
	}
	return base + "/" + strings.TrimLeft(key, "/")
}

func (r *r2Store) PresignPut(ctx context.Context, key, contentType string, expiry time.Duration) (string, error) {
	if r.presigner == nil {
		return "", fmt.Errorf("presigner not ready")
	}
	if expiry <= 0 {
		expiry = 20 * time.Minute
	}
	in := &s3.PutObjectInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(key),
	}
	if contentType != "" {
		in.ContentType = aws.String(contentType)
	}
	out, err := r.presigner.PresignPutObject(ctx, in, s3.WithPresignExpires(expiry))
	if err != nil {
		return "", err
	}
	return out.URL, nil
}

func (r *r2Store) Copy(ctx context.Context, srcKey, dstKey, contentType string) error {
	copySource := r.bucket + "/" + srcKey
	in := &s3.CopyObjectInput{
		Bucket:            aws.String(r.bucket),
		Key:               aws.String(dstKey),
		CopySource:        aws.String(copySource),
		MetadataDirective: types.MetadataDirectiveReplace,
		CacheControl:      aws.String("public, max-age=31536000, immutable"),
	}
	if contentType != "" {
		in.ContentType = aws.String(contentType)
	}
	_, err := r.client.CopyObject(ctx, in)
	return err
}

func (r *r2Store) Head(ctx context.Context, key string) (int64, bool, error) {
	out, err := r.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return 0, false, nil
	}
	var size int64
	if out.ContentLength != nil {
		size = *out.ContentLength
	}
	return size, true, nil
}

var mediaObjectStore ObjectStore = noopStore{}

func initObjectStore() {
	backend := strings.ToLower(strings.TrimSpace(os.Getenv("SURFING_OBJECT_STORE")))
	if backend == "" || backend == "local" || backend == "pvc" {
		mediaObjectStore = noopStore{}
		log.Info("ObjectStore: local PVC (SURFING_OBJECT_STORE unset/local)")
		return
	}
	if backend != "r2" {
		log.Warnf("ObjectStore: unknown SURFING_OBJECT_STORE=%q — using local PVC", backend)
		mediaObjectStore = noopStore{}
		return
	}

	accountID := strings.TrimSpace(os.Getenv("R2_ACCOUNT_ID"))
	accessKey := strings.TrimSpace(firstNonEmpty(os.Getenv("R2_ACCESS_KEY_ID"), os.Getenv("AWS_ACCESS_KEY_ID")))
	secretKey := strings.TrimSpace(firstNonEmpty(os.Getenv("R2_SECRET_ACCESS_KEY"), os.Getenv("AWS_SECRET_ACCESS_KEY")))
	bucket := strings.TrimSpace(firstNonEmpty(os.Getenv("R2_BUCKET"), os.Getenv("R2_BUCKET_NAME")))
	publicBase := strings.TrimSpace(firstNonEmpty(os.Getenv("R2_PUBLIC_BASE_URL"), os.Getenv("SURFING_CDN_BASE_URL")))
	endpoint := strings.TrimSpace(os.Getenv("R2_ENDPOINT"))

	if endpoint == "" && accountID != "" {
		endpoint = fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountID)
	}
	if s3URL := strings.TrimSpace(os.Getenv("R2_S3_URL")); s3URL != "" {
		if u, err := url.Parse(s3URL); err == nil && u.Host != "" {
			endpoint = u.Scheme + "://" + u.Host
			if bucket == "" {
				bucket = strings.Trim(u.Path, "/")
				if i := strings.Index(bucket, "/"); i >= 0 {
					bucket = bucket[:i]
				}
			}
		}
	}

	if endpoint == "" || accessKey == "" || secretKey == "" || bucket == "" {
		log.Warn("ObjectStore: R2 requested but missing endpoint/keys/bucket — falling back to local PVC")
		mediaObjectStore = noopStore{}
		return
	}
	if publicBase == "" {
		log.Warn("ObjectStore: R2_PUBLIC_BASE_URL unset — uploads will still write R2 but URLs stay /serve until set")
	}

	client := s3.New(s3.Options{
		Region:       firstNonEmpty(os.Getenv("R2_REGION"), "auto"),
		Credentials:  credentials.NewStaticCredentialsProvider(accessKey, secretKey, ""),
		BaseEndpoint: aws.String(endpoint),
		UsePathStyle: true,
	})

	mediaObjectStore = &r2Store{
		client:     client,
		presigner:  s3.NewPresignClient(client),
		bucket:     bucket,
		publicBase: publicBase,
	}
	log.Infof("ObjectStore: R2 ready bucket=%s endpoint=%s public=%s (presign enabled)", bucket, endpoint, publicBase)
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if strings.TrimSpace(v) != "" {
			return strings.TrimSpace(v)
		}
	}
	return ""
}

func mediaObjectKey(dayID, mediaID, ext string) string {
	return mediaObjectKeyPrefixed(dayID, mediaID, ext, "original")
}

func mediaDraftObjectKey(dayID, mediaID, ext string) string {
	return mediaObjectKeyPrefixed(dayID, mediaID, ext, "draft")
}

func mediaObjectKeyPrefixed(dayID, mediaID, ext, folder string) string {
	ext = strings.ToLower(ext)
	if ext != "" && !strings.HasPrefix(ext, ".") {
		ext = "." + ext
	}
	dayID = strings.TrimSpace(dayID)
	if dayID == "" {
		dayID = "_unscoped"
	}
	if folder == "" {
		folder = "original"
	}
	return "surfing/albums/" + dayID + "/" + folder + "/" + mediaID + ext
}

func putMediaObject(dayID, mediaID, ext, contentType string, data []byte) (objectKey, publicURL string, err error) {
	if !mediaObjectStore.Enabled() {
		return "", "", nil
	}
	key := mediaObjectKey(dayID, mediaID, ext)
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	if err := mediaObjectStore.Put(ctx, key, data, contentType); err != nil {
		return "", "", err
	}
	return key, mediaObjectStore.PublicURL(key), nil
}

func deleteMediaObject(dayID, mediaID, ext string) {
	if !mediaObjectStore.Enabled() {
		return
	}
	ext = normalizeExt(ext)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	for _, key := range []string{
		mediaObjectKey(dayID, mediaID, ext),
		mediaDraftObjectKey(dayID, mediaID, ext),
		"media/" + mediaID + ext,
	} {
		if err := mediaObjectStore.Delete(ctx, key); err != nil {
			log.Warnf("ObjectStore: delete %s failed: %v", key, err)
		}
	}
}

func promoteDraftObject(dayID, mediaID, ext, contentType string) (objectKey, publicURL string, err error) {
	if !mediaObjectStore.Enabled() {
		return "", "", fmt.Errorf("object store disabled")
	}
	ext = normalizeExt(ext)
	src := mediaDraftObjectKey(dayID, mediaID, ext)
	dst := mediaObjectKey(dayID, mediaID, ext)
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	if err := mediaObjectStore.Copy(ctx, src, dst, contentType); err != nil {
		return "", "", err
	}
	_ = mediaObjectStore.Delete(ctx, src)
	return dst, mediaObjectStore.PublicURL(dst), nil
}
