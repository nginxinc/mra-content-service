#!/bin/sh
NGINX_PID="/var/run/nginx.pid"    # /   (root directory)
NGINX_CONF="";
APP="/usr/local/go/bin/go run main.go error.go handlers.go router.go routes.go album_manager.go"

if [ ! -f .env ]; then
    touch .env
fi

su content-service -c "$APP" &

sleep 10
#APP gets rendered as go
APP_PID=`ps aux | grep "$APP" | grep -v grep`

case "$NETWORK" in
    fabric)
        NGINX_CONF="/etc/nginx/fabric_nginx_$CONTAINER_ENGINE.conf"
        echo 'Fabric configuration set'
        nginx -c "$NGINX_CONF" -g "pid $NGINX_PID;" &

        sleep 30

        while [ -f "$NGINX_PID" ] &&  [ "$APP_PID" ];
        do
	        sleep 5;
	        APP_PID=`ps aux | grep "$APP" | grep -v grep`;
        done
        ;;
    router-mesh)
        sleep 30

        while [ "$APP_PID" ];
        do
	        sleep 5;
	        APP_PID=`ps aux | grep "$APP" | grep -v grep`;
        done
        ;;
    *)
        echo 'Network not supported'
        exit 1
esac
