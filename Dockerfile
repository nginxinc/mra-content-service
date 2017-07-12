FROM golang:1.8-onbuild

ENV USE_NGINX_PLUS true


# Get other files required for installation
RUN apt-get update && apt-get install -y \
    wget \
    curl \
    apt-transport-https \
    git \
    libcurl3-gnutls \
    lsb-release \
    unzip \
    ca-certificates

# Install vault client
RUN wget -q https://releases.hashicorp.com/vault/0.6.0/vault_0.6.0_linux_amd64.zip && \
	  unzip -d /usr/local/bin vault_0.6.0_linux_amd64.zip

# Download certificate and key from the the vault and copy to the build context
ENV VAULT_TOKEN=4b9f8249-538a-d75a-e6d3-69f5355c1751 \
    VAULT_ADDR=http://vault.mra.nginxps.com:8200

RUN mkdir -p /etc/ssl/nginx && \
	vault token-renew && \
	vault read -field=value secret/nginx-repo.crt > /etc/ssl/nginx/nginx-repo.crt && \
	vault read -field=value secret/nginx-repo.key > /etc/ssl/nginx/nginx-repo.key && \
	vault read -field=value secret/ssl/csr.pem > /etc/ssl/nginx/csr.pem && \
	vault read -field=value secret/ssl/certificate.pem > /etc/ssl/nginx/certificate.pem && \
	vault read -field=value secret/ssl/key.pem > /etc/ssl/nginx/key.pem && \
	vault read -field=value secret/ssl/dhparam.pem > /etc/ssl/nginx/dhparam.pem

# Install nginx
ADD install-nginx.sh /usr/local/bin/
RUN /usr/local/bin/install-nginx.sh

RUN chown -R nginx /var/log/nginx/

# forward request logs to Docker log collector
RUN ln -sf /dev/stdout /var/log/nginx/access.log && \
    ln -sf /dev/stderr /var/log/nginx/error.log

# install application
# TODO
#COPY src/ /go/

#RUN ln -sf /dev/stdout /ingenious-pages/app/logs/prod.log && \
#    chown -R nginx:www-data /ingenious-pages/ && \
#    chmod -R 775 /ingenious-pages

#RUN cd /ingenious-pages && php composer.phar install --no-dev --optimize-autoloader

#COPY php5-fpm.conf /etc/php5/fpm/php-fpm.conf
#COPY php.ini /usr/local/etc/php/
COPY nginx /etc/nginx/

# Install and run NGINX config generator
RUN wget -q https://s3-us-west-1.amazonaws.com/fabric-model/config-generator/generate_config
RUN chmod +x generate_config && \
    ./generate_config -p /etc/nginx/fabric_config.yaml -t /etc/nginx/nginx-fabric.conf.j2 > /etc/nginx/nginx-fabric.conf

CMD ["./start.sh"]

EXPOSE 80 443