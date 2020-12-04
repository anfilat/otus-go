#!/usr/bin/env bash

# Если сокет закрылся со стороны сервера, то при следующей попытке отправить сообщение программа должна завершаться
# (допускается завершать программу после "неудачной" отправки нескольких сообщений).

set -xeuo pipefail

go build -o go-telnet

(echo -e "Hello\nFrom\nNC\n" && cat 2>/dev/null) | nc -l localhost 4242 >/tmp/nc.out &
NC_PID=$!

sleep 1
(echo -e "I\nam\nTELNET client\n" && cat 2>/dev/null) | ./go-telnet --timeout=5s localhost 4242 >/tmp/telnet.out &
TL_PID=$!

sleep 5
kill ${NC_PID} 2>/dev/null || true
echo "next message" > /proc/$TL_PID/fd/0

sleep 1
if ps -p $TL_PID > /dev/null
then
  echo "telnet client must be stopped"
  exit 1
fi

rm -f go-telnet
echo "PASS"
