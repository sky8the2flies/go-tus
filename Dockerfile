FROM golang:1.13

RUN mkdir -p /go/src/github.com/sky8the2flies/go-tus

WORKDIR /go/src/github.com/sky8the2flies/go-tus

RUN go get github.com/stretchr/testify
RUN go get github.com/tus/tusd
RUN go get github.com/syndtr/goleveldb/leveldb
