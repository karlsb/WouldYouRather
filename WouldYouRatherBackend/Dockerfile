# Base container on a golang image.
FROM golang:1.23 AS builder


# set directory inside container
WORKDIR /app


COPY go.mod go.sum ./
RUN go mod download



COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# second stage create machine
FROM ubuntu:24.04

WORKDIR /app
COPY --from=builder /app/main /app/main
# Seed database; copied onto the Fly volume on first boot (see seeddb.go)
COPY --from=builder /app/build-database/wouldyourather.db /app/build-database/

EXPOSE 8080

CMD ["./main"]
