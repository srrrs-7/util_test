FROM golang:latest AS builder

ADD ./concurrency/src /go/src
WORKDIR /go/src

ARG GOOS=linux
ARG GOARCH=amd64
ARG CGO_ENABLED=false
RUN go build -tags=local -o /go/bin/api /go/src/cmd/api

FROM scratch
COPY --from=builder /go/bin/api /usr/local/bin/api

EXPOSE 8080

CMD [ "/usr/local/bin/api" ]