FROM golang:1.27 AS build

WORKDIR /app

COPY main.go .
COPY go.mod .

RUN CGO_ENABLED=0 GOOS=linux go build -o app .

FROM alpine:3.24

WORKDIR /app

COPY --from=build /app/app .

CMD ["./app"]
