#!/usr/bin/env bash

# с нулевым таймаутом клиент не должен коннектиться к серверу

set -xeuo pipefail

go build -o go-telnet

(echo -e "Hello\nFrom\nNC\n" && cat 2>/dev/null) | nc -l localhost 4242 >/tmp/nc.out &
NC_PID=$!

sleep 1
(echo -e "I\nam\nTELNET client\n" && cat 2>/dev/null) | ./go-telnet --timeout=0s localhost 4242 2>/tmp/telnet.out &
TL_PID=$!

sleep 5
if ps -p $TL_PID > /dev/null
then
  echo "telnet client must be stopped"
  exit 1
fi
kill ${NC_PID} 2>/dev/null || true

function fileEquals() {
  local fileData
  fileData=$(cat "$1")
  [ "${fileData}" = "${2}" ] || (echo -e "unexpected output, $1:\n${fileData}" && exit 1)
}

fileEquals /tmp/telnet.out "cannot connect: dial tcp: i/o timeout"

rm -f go-telnet
echo "PASS"
