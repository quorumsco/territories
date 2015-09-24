FROM golang
MAINTAINER Wilmot Guillaume - Quorums

RUN go get github.com/tools/godep

ADD . /go/src/github.com/quorumsco/contacts

WORKDIR /go/src/github.com/quorumsco/contacts

RUN godep go build

EXPOSE 8080

ENTRYPOINT ["./contacts"]
