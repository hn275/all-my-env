# syntax=docker/dockerfile:1

FROM golang:1.20.5-alpine3.18
WORKDIR app
COPY . .
RUN go install github.com/cosmtrek/air@latest
RUN go get
CMD ["air"]
