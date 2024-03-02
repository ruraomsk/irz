#!/bin/bash
echo 'Compiling for Omsk teleofis'
GOOS=linux GOARCH=arm  go build
if [ $? -ne 0 ]; then
	echo 'An error has occurred! Aborting the script execution...'
	exit 1
fi
echo 'Copy irz to device'
scp irz root@192.168.115.26:/root
scp goirz.sh root@192.168.115.26:/root
# scp rc.local root@192.168.88.1:/etc
