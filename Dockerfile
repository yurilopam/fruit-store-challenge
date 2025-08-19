FROM golang:1.24.6-alpine

WORKDIR /app

RUN apk add --no-cache gcc musl-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV CGO_ENABLED=1
