FROM golang:1.7.4-alpine

RUN mkdir -p $GOPATH/src/github.com/Quynn-Software-Solutions/FOSSBounce

WORKDIR "$GOPATH/src/github.com/Quynn-Software-Solutions/FOSSBounce"

COPY . .

RUN go get

ENTRYPOINT ["/go/bin/FOSSBounce"]

EXPOSE 8080