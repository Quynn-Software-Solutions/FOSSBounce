FROM golang:1.22-alpine

RUN mkdir -p $GOPATH/src/github.com/Quynn-Software-Solutions/FOSSBounce/test

WORKDIR "$GOPATH/src/github.com/Quynn-Software-Solutions/FOSSBounce/test"

COPY . .

RUN go get

ENTRYPOINT ["/go/bin/FOSSBounce"]

EXPOSE 8080