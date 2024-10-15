# Stage 1: Build the Go program
FROM golang:1.19 AS build

WORKDIR /app

# Copy the Go source code to the container
COPY main.go .
COPY go.mod .

# Build the Go program
RUN CGO_ENABLED=0 GOOS=linux go build -o app .

# Stage 2: Copy the compiled binary to an Alpine Linux image
FROM alpine:3.20.3

WORKDIR /app

# Copy the compiled binary from the previous stage
COPY --from=build /app/app .

# Run the Go program
CMD ["./app"]
