#!/bin/bash

docker run -it --rm \
  --user 1000 \
  -v $PWD:/build znly/protoc \
  -I /build/ \
  --go_out=/build/ \
  /build/message.proto

mv message.pb.go ../go/
