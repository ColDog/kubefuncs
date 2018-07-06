#!/bin/bash


commit=$(git rev-parse --short HEAD)
tag="localhost:5000/example:$commit"

docker build -t $tag .

helm upgrade \
  --install \

