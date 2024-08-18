FROM golang:latest AS builder

ADD ./goapi/src /go/src
WORKDIR /go/src

ARG GOARCH='arm64'
ARG GOOS='linux'
RUN go build -ldflags="-s -w" -gcflags="-N -l" -buildmode="pie" \
    -o /go/bin/api /go/src/cmd/api


FROM scratch
COPY --from=builder /go/bin/api /usr/local/bin/api

EXPOSE 8080

CMD [ "/usr/local/bin/api" ]