# syntax=docker/dockerfile:1

FROM rust:1-buster
WORKDIR /app
COPY . .
RUN cargo build --release
CMD ["./target/release/server"]
