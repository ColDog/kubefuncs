#!/bin/bash

ip=$(minikube ip | tr -d '\n')

echo "sending request"
curl --fail -i -H 'Host: gateway.local' http://$ip/test/hello || exit 1

wrk -H 'Host: gateway.local' http://$ip/test/hello
