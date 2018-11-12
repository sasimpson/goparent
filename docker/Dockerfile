FROM golang:alpine as builder

WORKDIR /go/src/github.com/sasimpson/goparent

COPY . .
ENV GO111MODULE=on
RUN go build -o goparent-service /go/src/github.com/sasimpson/goparent/cmd/goparent-service/main.go

FROM alpine:latest

COPY --from=builder /go/src/github.com/sasimpson/goparent/goparent-service .
COPY goparent_sample.json /etc/config/goparent.json

ENTRYPOINT [ "./goparent-service" ]
EXPOSE 8000