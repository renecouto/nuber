#!/bin/bash

function fetchroute() {
    curl -sS -XPOST -H 'Content-Type:application/json' -H 'USER-TOKEN:123' localhost:8080/routes/_fetch -d '{
    "Origin": {
        "X": 10,
        "Y": 10
    },
    "Destination": {
        "X": 10,
        "Y": 10
    }
}'

}

ROUTEINFO=$(fetchroute)
echo "got routeinfo: $ROUTEINFO"



function askride() {
    curl -XPOST -H 'Content-Type:application/json' -H 'USER-TOKEN:123' localhost:8080/rides/_ask -d "$ROUTEINFO"
}

RIDEIDINFO=$(askride)
echo GOT RIDE $RIDEIDINFO
echo "waiting a bit"
sleep 10

echo "trying to accept ride with driver id 64"
curl -sS -XPOST -H 'Content-Type:application/json' -H 'USER-TOKEN:64' localhost:8080/rides/_accept -d "$RIDEIDINFO"
echo "accepted ride with driver id 64"