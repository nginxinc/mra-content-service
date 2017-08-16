FROM golang:1.8.3-jessie

ENV USE_NGINX_PLUS=true \
    USE_VAULT=true \
    USE_LOCAL=false

RUN mkdir -p /go/src/app
WORKDIR /go/src/app

# this will ideally be built by the ONBUILD below ;)
CMD ["go-wrapper", "run"]

COPY app /go/src/app/
COPY nginx/ssl /etc/ssl/nginx/
COPY vault_env.sh /etc/letsencrypt/
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

# Install nginx
ADD install-nginx.sh /usr/local/bin/
COPY nginx /etc/nginx/
RUN /usr/local/bin/install-nginx.sh && \
# forward request logs to Docker log collector
    ln -sf /dev/stdout /var/log/nginx/access.log && \
    ln -sf /dev/stderr /var/log/nginx/error.log

#ADD app /app/

EXPOSE 80 443 12002

ENTRYPOINT ["./start.sh"]