FROM golang:1.19-alpine3.17 AS builder

RUN apk update && apk add --no-cache ca-certificates

FROM scratch

WORKDIR /bin

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY /bin/shop shop

ENTRYPOINT ["kitchen"]
