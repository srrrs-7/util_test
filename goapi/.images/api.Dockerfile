FROM golang:latest AS builder

ADD ./src /go/src
WORKDIR /go/src

ARG GOARCH='amd64'
ARG GOOS='linux'
RUN go build -ldflags="-s -w" -gcflags="-N -l" -buildmode="pie" \
    -o /go/bin/api /go/src/cmd/api


FROM scrach AS runner
COPY --from=builder /go/bin/api /usr/local/bin/api

CMD [ "/usr/local/bin/api" ]