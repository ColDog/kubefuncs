#!/bin/bash

POD=$(kubectl -n kubefuncs get pods | grep gateway | awk '{print $1}' | head -n 1)
kubectl -n kubefuncs port-forward "$POD" 8080:8080 &
pid=$!

sleep 1

echo "sending request"
curl --fail -i localhost:8080/test/hello || {
  kill $pid
  exit 1
}

wrk http://localhost:8080/test/hello

kill $pid
