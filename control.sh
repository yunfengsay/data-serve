#!/bin/sh

NAME="fileserveforshare"
BINARY="node $(dirname $_)/index.js --$NAME"
PIDFILE="logs/pid.txt"

stop() {
    pid=$(ps -ef |grep $NAME |grep -v grep | awk '{print $2}') 
    kill $pid
}

start() {
    stop
    eval $BINARY &
    echo "$NAME is runinig"
    exit 0
}
# printenv | more
dirname $_

$1