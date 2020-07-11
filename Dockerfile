FROM golang:rc-alpine AS builder
WORKDIR /go/src/github.com/tutilus/gotwarden
# Update from go.mod all the dependencies
# RUN go get -d -v golang.org/x/net/html  
# Copy all the source
RUN apk add --no-cache musl-dev git bash make gcc
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

WORKDIR /go/src/github.com/tutilus/gotwarden

ADD . .

# Download dependencies
RUN go mod download

# Build
RUN make build

FROM alpine:latest
ENV DB_FILEPATH /gotwarden/db/warden.db
ENV PORT 3000
ENV GIN_MODE release

RUN apk add --no-cache ca-certificates
RUN apk add --no-cache musl sqlite
WORKDIR /gotwarden/bin
COPY --from=builder /go/src/github.com/tutilus/gotwarden/gotwarden server
# Copy .env (version Prod)
RUN mkdir /gotwarden/db
VOLUME /gotwarden/db
ENTRYPOINT ["/gotwarden/bin/server"]