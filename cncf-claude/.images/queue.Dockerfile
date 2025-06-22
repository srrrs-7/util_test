FROM golang:1.24-alpine AS builder

COPY ./src /go/src
WORKDIR /go/src

ARG GOOS=linux
ARG GOARCH=amd64
ARG CGO_ENABLED=false
RUN go build -ldflags="-s -w" -o /go/bin/queue ./cmd/queue


# claude code container
FROM alpine:latest

COPY --from=builder /go/bin/queue /usr/local/bin/queue

CMD [ "/usr/local/bin/queue" ]
