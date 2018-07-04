FROM golang:1.10-alpine as builder
COPY gateway /go/src/github.com/coldog/kubefuncs/gateway
COPY clients /go/src/github.com/coldog/kubefuncs/clients
COPY vendor /go/src/github.com/coldog/kubefuncs/vendor
RUN go build \
  -o /build/gateway \
  github.com/coldog/kubefuncs/gateway

FROM alpine:3.7
RUN apk add --no-cache ca-certificates
COPY --from=builder /build/gateway /bin/gateway
