FROM golang:alpine as builder

WORKDIR /go/src/github.com/sasimpson/goparent

RUN apk add git curl --no-cache
RUN curl -s https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

COPY . .
RUN dep ensure
RUN go build -o goparent-service ./cmd/goparent-service/main.go

FROM alpine:latest

COPY --from=builder /go/src/github.com/sasimpson/goparent/goparent-service .
COPY goparent_sample.json /etc/config/goparent.json

ENTRYPOINT [ "./goparent-service" ]
EXPOSE 8000