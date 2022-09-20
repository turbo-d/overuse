##
## Build
##

FROM golang:1.14.6-alpine3.12 as builder

WORKDIR /go/src/github.com/turbo-d/overuse

COPY go.mod go.sum /go/src/github.com/turbo-d/overuse/

RUN go mod download

COPY . /go/src/github.com/turbo-d/overuse

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o build/overuse github.com/turbo-d/overuse

##
## Deploy
##

FROM alpine

RUN apk add --no-cache ca-certificates && update-ca-certificates

COPY --from=builder /go/src/github.com/turbo-d/overuse/build/overuse /usr/bin/overuse

EXPOSE 8080 8080

ENTRYPOINT ["/usr/bin/overuse"]
