# surfing-r2 is created out-of-band (never commit keys):
#   oc -n surfing-service-system create secret generic surfing-r2 \
#     --from-env-file=/home/dasm/r2_creds_dasmlab_surfing
#
# Required keys: R2_ACCOUNT_ID R2_ACCESS_KEY_ID R2_SECRET_ACCESS_KEY
#                R2_BUCKET R2_PUBLIC_BASE_URL  (+ optional R2_S3_URL)
