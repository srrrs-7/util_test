FROM golang:1.24 AS builder

COPY ./src /go/src
WORKDIR /go/src
ARG GOOS=linux
ARG GOARCH=amd64
ARG CGO_ENABLED=false
RUN go build --ldflags '-extldflags "-static"' -o /go/bin/proxy ./cmd/proxy

FROM alpine:3
COPY --from=builder /go/bin/proxy /usr/local/bin/proxy

CMD [ "/usr/local/bin/proxy" ]