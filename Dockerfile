FROM golang:1.17.8-alpine3.14  as build-env

ENV GO111MODULE=on
ENV BUILDPATH=github.com/icowan/shorter
ENV GOPROXY=https://goproxy.cn
ENV GOPATH=/go
RUN mkdir -p /go/src/${BUILDPATH}
COPY ./ /go/src/${BUILDPATH}
RUN cd /go/src/${BUILDPATH} && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go install ./cmd/

FROM alpine:3.14.6

RUN apk add --no-cache \
		ca-certificates \
		curl \
		tzdata \
		&& cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
		&& echo "Asia/Shanghai" > /etc/timezone \
		&& apk del tzdata \
		&& rm -rf /var/cache/apk/*

COPY --from=build-env /go/bin/cmd /go/bin/shorter
COPY ./dist /go/bin/dist

WORKDIR /go/bin/
CMD ["/go/bin/shorter", "-http-addr", ":8080"]