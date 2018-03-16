FROM golang:alpine

WORKDIR /go/src/github.com/sasimpson/goparent
COPY . .
COPY goparent_sample.json /etc/config/goparent.json

RUN apk update && apk add git curl --no-cache
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
RUN dep ensure
RUN go install -v ./...

ENTRYPOINT [ "/go/bin/goparent-service" ]

EXPOSE 8000
