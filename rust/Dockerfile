# Build stage
FROM rust:1-bookworm AS builder

WORKDIR /app
COPY . .

RUN cargo build --release

# Run stage
FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    libssl3 \
    wget \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY --from=builder /app/target/release/rust-microservice .
EXPOSE 8081

CMD ["./rust-microservice"]
