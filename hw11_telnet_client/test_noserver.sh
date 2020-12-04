#!/usr/bin/env bash

# При подключении к несуществующему серверу, программа должна завершаться с ошибкой соединения/таймаута.

set -xeuo pipefail

go build -o go-telnet

(echo -e "I\nam\nTELNET client\n" && cat 2>/dev/null) | ./go-telnet --timeout=3s "127.0.0.1" 4242 2>/tmp/telnet.out &
TL_PID=$!

sleep 5
if ps -p $TL_PID > /dev/null
then
  echo "telnet client must be stopped"
  exit 1
fi

function fileEquals() {
  local fileData
  fileData=$(cat "$1")
  [ "${fileData}" = "${2}" ] || (echo -e "unexpected output, $1:\n${fileData}" && exit 1)
}

fileEquals /tmp/telnet.out "cannot connect: dial tcp 127.0.0.1:4242: connect: connection refused"

rm -f go-telnet
echo "PASS"
