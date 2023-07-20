# First stage downloads majority of dependencies
# It is rarely changed during development
# Runs every time go.mod or go.sum are changed.
FROM golang:1.20-alpine as build_dependencies
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

# The second stage builds application.
# It also downloads the rest of the dependencies that can be more dynamic.
# Runs every time the code is changed.
FROM golang:1.20-alpine as build
WORKDIR /app
COPY --from=build_dependencies /go /go

COPY ./internal ./internal
COPY ./cmd ./cmd
COPY go.mod go.sum ./
RUN go mod tidy

RUN go build -o main ./cmd/main.go

# Runnable stage. Have minimum dependencies and files.
FROM alpine:3.14
WORKDIR /app
RUN apk update
RUN apk add --no-cache bash

COPY --from=build /app/main ./

CMD [ "/app/main" ]
