FROM golang:1.10-alpine as builder
COPY example /go/src/github.com/coldog/kubefuncs/example
COPY clients /go/src/github.com/coldog/kubefuncs/clients
COPY vendor /go/src/github.com/coldog/kubefuncs/vendor
RUN go build \
  -o /build/http-pong \
  github.com/coldog/kubefuncs/example/http-pong
RUN go build \
  -o /build/http-archive \
  github.com/coldog/kubefuncs/example/http-archive

FROM alpine:3.7
RUN apk add --no-cache ca-certificates
COPY --from=builder /build/http-pong /bin/http-pong
COPY --from=builder /build/http-archive /bin/http-archive
