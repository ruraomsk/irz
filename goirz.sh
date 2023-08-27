#!/bin/ash
while true
do
    echo "start controller" >> start
    ./irz > /dev/null 2>/dev/null
    echo "need restart " >> start
done 
