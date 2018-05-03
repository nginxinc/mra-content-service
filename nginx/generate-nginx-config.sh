#!/usr/bin/env bash

wget -O /usr/local/sbin/generate_config -q https://s3-us-west-1.amazonaws.com/fabric-model/config-generator/generate_config
chmod +x /usr/local/sbin/generate_config

CONFIG_FILE=/etc/nginx/fabric/fabric_config.yaml

echo -e "\033[32m -----"
echo -e "\033[32m Building for ${CONTAINER_ENGINE}"
echo -e "\033[32m -----\033[0m"

case "$CONTAINER_ENGINE" in
    kubernetes)
        CONFIG_FILE=/etc/nginx/fabric/fabric_config_k8s.yaml
        ;;
    local)
        CONFIG_FILE=/etc/nginx/fabric/fabric_config_local.yaml
        ;;
esac

/usr/local/sbin/generate_config -p ${CONFIG_FILE} -t /etc/nginx/fabric/fabric_nginx.conf.j2 > /etc/nginx/nginx.conf
