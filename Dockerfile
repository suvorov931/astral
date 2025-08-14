FROM golang:1.24-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o astral cmd/*.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app/

COPY --from=build /app/astral .
COPY config/config.env /app/config/config.env
COPY database/migrations /app/database/migrations

CMD ["./astral"]