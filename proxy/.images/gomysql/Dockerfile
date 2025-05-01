FROM golang:1.24 AS builder

COPY ./src /go/src
WORKDIR /go/src
ARG GOOS=linux
ARG GOARCH=amd64
ARG CGO_ENABLED=false
RUN go build --ldflags '-extldflags "-static"' -o /go/bin/gomysql ./cmd/gomysql

FROM alpine:3
COPY --from=builder /go/bin/gomysql /usr/local/bin/gomysql

CMD [ "/usr/local/bin/gomysql" ]