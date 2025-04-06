FROM golang:latest AS builder

WORKDIR /app

# Install gqlgen
RUN go install github.com/99designs/gqlgen@latest

ENTRYPOINT ["gqlgen"]
