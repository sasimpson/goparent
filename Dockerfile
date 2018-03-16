FROM golang:alpine

WORKDIR /go/src/github.com/sasimpson/goparent
COPY . .
COPY goparent_sample.json /etc/config/goparent.json

RUN go install -v ./...

ENTRYPOINT [ "/go/bin/goparent-service" ]

EXPOSE 8000
