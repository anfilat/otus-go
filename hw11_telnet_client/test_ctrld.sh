#!/usr/bin/env bash

#При нажатии Ctrl+D программа должна закрывать сокет и завершаться с сообщением.

set -xeuo pipefail

go build -o go-telnet

nc -l localhost 4242 >/tmp/nc.out &
NC_PID=$!

sleep 1
echo -e "I am TELNET client\n" >/tmp/input.in
./go-telnet --timeout=5s localhost 4242 2>/tmp/telnet.err </tmp/input.in &
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

fileEquals /tmp/nc.out "I am TELNET client"

expected_telnet_err='...Connected to localhost:4242
...EOF'
fileEquals /tmp/telnet.err "${expected_telnet_err}"

rm -f go-telnet
echo "PASS"
