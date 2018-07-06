FROM golang:1.10-alpine as builder
COPY example /go/src/github.com/coldog/kubefuncs/example
COPY clients /go/src/github.com/coldog/kubefuncs/clients
COPY vendor /go/src/github.com/coldog/kubefuncs/vendor
RUN go build \
  -o /build/app \
  github.com/coldog/kubefuncs/example

FROM alpine:3.7
RUN apk add --no-cache ca-certificates
COPY --from=builder /build/app /bin/app
