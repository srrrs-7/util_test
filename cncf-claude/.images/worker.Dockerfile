FROM golang:1.24-alpine AS builder

COPY ./src /go/src
WORKDIR /go/src

ARG GOOS=linux
ARG GOARCH=amd64
ARG CGO_ENABLED=false
RUN go build -ldflags="-s -w" -o /go/bin/worker ./cmd/worker


# claude code container
FROM oven/bun

RUN apt update && apt install -y nodejs
RUN bun install -g @anthropic-ai/claude-code

COPY --from=builder /go/bin/worker /usr/local/bin/worker

CMD [ "/usr/local/bin/worker" ]
