FROM golang:1.9.2

WORKDIR /go/src/github.com/sasimpson/goparent
COPY . .
COPY goparent_sample.json /etc/config/goparent.json

RUN go get -d -v ./...
RUN go install -v ./...

CMD [ "ls -l /go/bin/" ]
ENTRYPOINT [ "/go/bin/goparent" ]

EXPOSE 8000
