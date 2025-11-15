FROM node:22 AS build
# Build-time args (populated by GitHub Actions)
ARG COUNTER_API_TOKEN
ARG COUNTER_WORKSPACE
ARG COUNTER_COUNTER_SLUG
ARG COUNTER_FULL_PATH
ARG APP_VERSION
WORKDIR /app
COPY . .
# Ensure envs are available during build for quasar build-time injection
ENV COUNTER_API_TOKEN=${COUNTER_API_TOKEN}
ENV COUNTER_WORKSPACE=${COUNTER_WORKSPACE}
ENV COUNTER_COUNTER_SLUG=${COUNTER_COUNTER_SLUG}
ENV COUNTER_FULL_PATH=${COUNTER_FULL_PATH}
ENV APP_VERSION=${APP_VERSION}
RUN npm install && npm install counterapi && npm run build

FROM docker.io/library/nginx:alpine
# Propagate envs to runtime container (for future diagnostics or dynamic usage)
ARG COUNTER_API_TOKEN
ARG COUNTER_WORKSPACE
ARG COUNTER_COUNTER_SLUG
ARG COUNTER_FULL_PATH
ARG APP_VERSION
ENV COUNTER_API_TOKEN=${COUNTER_API_TOKEN}
ENV COUNTER_WORKSPACE=${COUNTER_WORKSPACE}
ENV COUNTER_COUNTER_SLUG=${COUNTER_COUNTER_SLUG}
ENV COUNTER_FULL_PATH=${COUNTER_FULL_PATH}
ENV APP_VERSION=${APP_VERSION}
# Configure nginx to run as root (for OpenShift anyuid SCC)
# Create temp directories and set permissions before nginx tries to chown them
RUN sed -i 's/user nginx;/user root;/' /etc/nginx/nginx.conf && \
    mkdir -p /var/cache/nginx/client_temp /var/cache/nginx/proxy_temp /var/cache/nginx/fastcgi_temp /var/cache/nginx/uwsgi_temp /var/cache/nginx/scgi_temp /var/run && \
    chmod -R 777 /var/cache/nginx /var/run && \
    chown -R root:root /var/cache/nginx /var/run || true
COPY --from=build /app/dist/spa /usr/share/nginx/html
EXPOSE 80

