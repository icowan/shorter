FROM golang:1.13.4-alpine3.10 as build-env

ENV GO111MODULE=on
ENV BUILDPATH=github.com/icowan/shorter
ENV GOPROXY=https://goproxy.cn
ENV GOPATH=/go
RUN mkdir -p /go/src/${BUILDPATH}
COPY ./ /go/src/${BUILDPATH}
RUN cd /go/src/${BUILDPATH} && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go install ./cmd/

FROM alpine:latest

RUN apk update \
        && apk upgrade \
        && apk add --no-cache \
        ca-certificates \
        curl \
        && update-ca-certificates 2>/dev/null || true

COPY --from=build-env /go/bin/cmd /go/bin/shorter

WORKDIR /go/bin/
CMD ["/go/bin/shorter", "-http-addr", ":8080"]