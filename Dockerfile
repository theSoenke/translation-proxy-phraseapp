FROM golang:1.10-alpine as builder

RUN apk update && apk upgrade
RUN apk add git

RUN mkdir -p /go/src/github.com/thesoenke/translation-proxy
COPY . /go/src/github.com/thesoenke/translation-proxy

WORKDIR /go/src/github.com/thesoenke/translation-proxy
RUN go get
RUN go build -o app


FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=builder /go/src/github.com/thesoenke/translation-proxy/app .
ENV GIN_MODE=release

CMD ["./app"]
