FROM golang:1.24 AS builder

COPY ./src /go/src
WORKDIR /go/src
ARG GOOS=linux
ARG GOARCH=amd64
ARG CGO_ENABLED=false
RUN go build --ldflags '-extldflags "-static"' -o /go/bin/client ./cmd/client

FROM alpine:3
COPY --from=builder /go/bin/client /usr/local/bin/client

RUN apk update && apk add mysql-client

CMD [ "/usr/local/bin/client" ]