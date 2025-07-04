FROM golang:latest AS builder

ADD ./src /go/src
WORKDIR /go/src

ARG GOOS=linux
ARG GOARCH=amd64
ARG CGO_ENABLED=false
RUN go build -tags=local -o /go/bin/client /go/src/cmd/client

FROM scratch
COPY --from=builder /go/bin/client /usr/local/bin/client

EXPOSE 8080

CMD [ "/usr/local/bin/client" ]