# Stage 1: Build
FROM docker.io/library/golang:latest AS builder

WORKDIR /workspace

ARG ARCH=amd64
ARG goproxy=https://proxy.golang.org
ENV GOPROXY=${goproxy}

COPY main-app/go.mod main-app/go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY main-app/ .

RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux GOARCH=${ARCH} \
    go build -ldflags "-s -w" \
    -o surfing-service

# Stage 2: Run
FROM docker.io/library/ubuntu:latest

RUN apt-get update && apt-get install -y ca-certificates curl && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY --from=builder /workspace/surfing-service .
RUN mkdir -p /data/media && chown -R 65532:65532 /app /data

USER 65532

ENV SURFING_DATA_DIR=/data \
    SURFING_API_PORT=10023
EXPOSE 10023 9222

ENTRYPOINT ["/app/surfing-service"]
