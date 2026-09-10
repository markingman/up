FROM golang:1.19 AS build

WORKDIR /app

COPY main.go .
COPY go.mod .

RUN CGO_ENABLED=0 GOOS=linux go build -o app .

FROM alpine:3.20.3

WORKDIR /app

COPY --from=build /app/app .

CMD ["./app"]
