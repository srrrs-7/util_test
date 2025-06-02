FROM golang:latest AS builder

ADD ./concurrency/src /go/src
WORKDIR /go/src

ARG GOOS=linux
ARG GOARCH=amd64
ARG CGO_ENABLED=false
RUN go build -tags=local -o /go/bin/worker /go/src/cmd/worker

FROM scratch
COPY --from=builder /go/bin/worker /usr/local/bin/worker

EXPOSE 8080

CMD [ "/usr/local/bin/worker" ]