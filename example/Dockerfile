FROM golang:1.10 as builder
COPY example /go/src/github.com/coldog/kubefuncs/example
COPY vendor /go/src/github.com/coldog/kubefuncs/vendor
RUN go build \
  -o /build/http-pong \
  github.com/coldog/kubefuncs/example/http-pong

FROM alpine:3.7
RUN apk add --no-cache ca-certificates
COPY --from=builder /build/http-pong /bin/http-pong
