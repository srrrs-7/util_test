FROM rust:slim-bookworm AS builder
RUN apt-get update && \
    apt-get install -y --no-install-recommends mold && \
    rm -rf /var/lib/apt/lists/*

ADD ./rustapi /rustapi
WORKDIR /rustapi/src

ENV RUSTFLAGS="-C link-arg=-fuse-ld=mold"
RUN --mount=type=cache,target=/usr/local/cargo/registry \
    --mount=type=cache,target=/app/target \
    cargo build --release --bin worker

FROM gcr.io/distroless/cc-debian12
COPY --from=builder /rustapi/src/target/release/worker /usr/local/bin/worker
EXPOSE 8080

CMD [ "/usr/local/bin/worker" ]