#!/bin/sh
NGINX_PID="/var/run/nginx.pid"    # /   (root directory)
APP="go run error.go handlers.go logger.go main.go router.go routes.go"

NGINX_CONF="/etc/nginx/nginx.conf";

if [ ! -f .env ]; then
    touch .env
fi

if [ "$NETWORK" = "fabric" ]
then
    echo fabric configuration set;
fi

$APP &

nginx -c "$NGINX_CONF" -g "pid $NGINX_PID;" &

sleep 10
#APP gets rendered as go
APP=go
APP_PID=`ps aux | grep "$APP" | grep -v grep`

./insert.sh

while [ -f "$NGINX_PID" ] &&  [ "$APP_PID" ];
do
	sleep 5;
	APP_PID=`ps aux | grep "$APP" | grep -v grep`;
done
