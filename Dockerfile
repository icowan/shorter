FROM golang:1.13.4-alpine3.10 as build-env

ENV GO111MODULE=on
ENV BUILDPATH=github.com/icowan/shortener
ENV GOPROXY=https://goproxy.cn
ENV GOPATH=/go
RUN mkdir -p /go/src/${BUILDPATH}
COPY ./ /go/src/${BUILDPATH}
RUN cd /go/src/${BUILDPATH} && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go install -v

FROM alpine:latest

RUN apk update \
        && apk upgrade \
        && apk add --no-cache \
        ca-certificates \
        curl \
        && update-ca-certificates 2>/dev/null || true

COPY --from=build-env /go/bin/shortener /go/bin/shortener

WORKDIR /go/bin/
CMD ["/go/bin/shortener", "start", "-p", ":8080", "-c", "/etc/shortener/app.cfg"]