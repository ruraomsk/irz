#!/bin/bash
echo 'Compiling'
GOOS=linux GOARCH=mipsle GOMIPS=softfloat go build
if [ $? -ne 0 ]; then
	echo 'An error has occurred! Aborting the script execution...'
	exit 1
fi
echo 'Copy irz to device'
scp irz root@192.168.2.1:/root
# echo 'Copy rc.local to device'
# scp rc.local root@192.168.2.1:/etc
#scp test.bin admin@192.168.115.29:/home/admin