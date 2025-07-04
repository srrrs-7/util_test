FROM golang:latest AS builder

ADD ./src /go/src
WORKDIR /go/src

ARG GOOS=linux
ARG GOARCH=amd64
ARG CGO_ENABLED=false
RUN go build -tags=local -o /go/bin/server /go/src/cmd/server

FROM scratch
COPY --from=builder /go/bin/server /usr/local/bin/server

EXPOSE 8080

CMD [ "/usr/local/bin/server" ]