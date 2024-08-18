FROM golang:latest AS builder

ADD ./goapi/src /go/src
WORKDIR /go/src

ARG GOOS=linux
ARG GOARCH=arm64
ARG CGO_ENABLED=false
RUN go build -ldflags "-s -w" -gcflags "-N" -buildmode "pie" \
    -o /go/bin/worker /go/src/cmd/worker

CMD ["/go/bin/worker"]


# FROM alpine:3
# COPY --from=builder /go/bin/worker /usr/local/bin/worker

# CMD [ "/usr/local/bin/worker" ]