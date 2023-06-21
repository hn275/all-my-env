# syntax=docker/dockerfile:1

FROM rust:1-buster
WORKDIR /app
COPY . .
RUN cargo install cargo-watch
CMD ["cargo", "watch", "-x", "run"]
