#!/bin/sh
APP="go run error.go handlers.go logger.go main.go router.go routes.go"

if [ ! -f .env ]; then
    touch .env
fi

if [ "$NETWORK" = "fabric" ]
then
    echo fabric configuration set;
    NGINX_PID="/var/run/nginx.pid"    # /   (root directory)
    NGINX_CONF="/etc/nginx/nginx.conf";
    nginx -c "$NGINX_CONF" -g "pid $NGINX_PID;" &
fi

$APP &

sleep 10
#APP gets rendered as go
APP=go
APP_PID=`ps aux | grep "$APP" | grep -v grep`

if [ "$NETWORK" = "fabric" ]
then
    while [ -f "$NGINX_PID" ] &&  [ "$APP_PID" ];
    do
	    sleep 5;
	    APP_PID=`ps aux | grep "$APP" | grep -v grep`;
    done
else
    while [ "$APP_PID" ];
    do
	    sleep 5;
	    APP_PID=`ps aux | grep "$APP" | grep -v grep`;
    done
fi

