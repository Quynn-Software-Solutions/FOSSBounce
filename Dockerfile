FROM golang:1.22-alpine

# Install git
RUN apk update && apk add --no-cache git

RUN mkdir -p $GOPATH/src/github.com/Quynn-Software-Solutions/FOSSBounce

WORKDIR "$GOPATH/src/github.com/Quynn-Software-Solutions/FOSSBounce"

COPY . .

RUN go get

# Build the Go application
RUN go build -o /go/bin/FOSSBounce .

ENTRYPOINT ["/go/bin/FOSSBounce"]

EXPOSE 8080