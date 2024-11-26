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
COPY .env /app/.env
COPY --from=builder /app/main /app/main
COPY --from=builder /app/build-database/*.db /app/build-database/

RUN ls -l /app
RUN chmod 666 /app/build-database/wouldyourather.db
RUN chmod 777 /app
RUN chmod 777 /app/build-database
# Expose the port your app runs on
EXPOSE 8080

# Start the Go application
CMD ["./main"]
