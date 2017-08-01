FROM golang:1.8-onbuild

ENV USE_NGINX_PLUS=true \
    VAULT_TOKEN=4b9f8249-538a-d75a-e6d3-69f5355c1751 \
    VAULT_ADDR=http://vault.mra.nginxps.com:8200

# Get other files required for installation
RUN apt-get update && apt-get install -y \
    wget \
    curl \
    apt-transport-https \
    git \
    libcurl3-gnutls \
    lsb-release \
    unzip \
    ca-certificates && \
# Install vault client
    wget -q https://releases.hashicorp.com/vault/0.6.0/vault_0.6.0_linux_amd64.zip && \
    unzip -d /usr/local/bin vault_0.6.0_linux_amd64.zip && \
    mkdir -p /etc/ssl/nginx 

# Install nginx
ADD install-nginx.sh /usr/local/bin/
COPY nginx /etc/nginx/
RUN /usr/local/bin/install-nginx.sh && \
# forward request logs to Docker log collector
    ln -sf /dev/stdout /var/log/nginx/access.log && \
    ln -sf /dev/stderr /var/log/nginx/error.log


CMD ["./start.sh"]

EXPOSE 80 443