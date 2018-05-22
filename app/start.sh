#!/bin/sh
NGINX_PID="/var/run/nginx.pid"    # /   (root directory)
NGINX_CONF="";
APP="/usr/local/go/bin/go run main.go error.go handlers.go logger.go router.go routes.go album_manager.go"

if [ ! -f .env ]; then
    touch .env
fi

case "$NETWORK" in
    fabric)
        NGINX_CONF="/etc/nginx/fabric_nginx_$CONTAINER_ENGINE.conf"
        echo 'Fabric configuration set'
        nginx -c "$NGINX_CONF" -g "pid $NGINX_PID;" &
        ;;
    router-mesh)
        ;;
    *)
        echo 'Network not supported'
esac

su content-service -c "$APP"

sleep 10
#APP gets rendered as go
APP_PID=`ps aux | grep "$APP" | grep -v grep`

while [ -f "$NGINX_PID" ] &&  [ "$APP_PID" ];
do
	sleep 5;
	APP_PID=`ps aux | grep "$APP" | grep -v grep`;
done
