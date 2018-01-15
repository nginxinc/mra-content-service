FROM golang:1.8.3-jessie

ARG CONTAINER_ENGINE_ARG
ARG USE_NGINX_PLUS_ARG
ARG USE_VAULT_ARG

# CONTAINER_ENGINE specifies the container engine to which the
# containers will be deployed. Valid values are:
# - kubernetes (default)
# - mesos
# - local
ENV USE_NGINX_PLUS=${USE_NGINX_PLUS_ARG:-true} \
    USE_VAULT=${USE_VAULT_ARG:-false} \
    CONTAINER_ENGINE=${CONTAINER_ENGINE_ARG:-kubernetes}

RUN mkdir -p /go/src/app && echo ${CONTAINER_ENGINE_ARG}
WORKDIR /go/src/app

# this will ideally be built by the ONBUILD below ;)
CMD ["go-wrapper", "run"]

COPY app /go/src/app/
COPY nginx/ssl /etc/ssl/nginx/

# Get other files required for installation
RUN go-wrapper download && \
    go-wrapper install && \
    apt-get update && apt-get install -y \
    apt-transport-https \
    ca-certificates \
    curl \
    git \
    libcurl3-gnutls \
    lsb-release \
    unzip \
    vim \
    wget && \
    mkdir -p /etc/ssl/nginx

# Install nginx and forward request logs to Docker log collector
ADD install-nginx.sh /usr/local/bin/
COPY nginx /etc/nginx/
RUN /usr/local/bin/install-nginx.sh && \
    ln -sf /dev/stdout /var/log/nginx/access_log && \
    ln -sf /dev/stderr /var/log/nginx/error_log

RUN ./test.sh

EXPOSE 80 443 12002

ENTRYPOINT ["./start.sh"]
