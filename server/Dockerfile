# syntax=docker/dockerfile:1

FROM golang:1.20.5-alpine3.18 as build
WORKDIR /app
COPY . .
RUN go build -o app main.go

FROM alpine:3.18 as prod
WORKDIR /app
COPY --from=build /app .
CMD ["./app"]
