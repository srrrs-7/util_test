FROM golang:latest AS builder

ADD ./src /go/src
WORKDIR /go/src

ARG GOARCH='amd64'
ARG GOOS='linux'
RUN go build -ldflags="-s -w" -gcflags="-N -l" -buildmode="pie" \
    -o /go/bin/worker /go/src/cmd/worker


FROM scrach AS runner
COPY --from=builder /go/bin/worker /usr/local/bin/worker

CMD [ "/usr/local/bin/worker" ]