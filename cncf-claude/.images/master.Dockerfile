FROM golang:1.24-alpine AS builder

COPY ./src /go/src
WORKDIR /go/src

ARG GOOS=linux
ARG GOARCH=amd64
ARG CGO_ENABLED=false
RUN go build -ldflags="-s -w" -o /go/bin/master ./cmd/master


# claude code container
FROM oven/bun

RUN apt update && apt install -y nodejs
RUN bun install -g @anthropic-ai/claude-code

COPY --from=builder /go/bin/master /usr/local/bin/master

CMD [ "/usr/local/bin/master" ]
