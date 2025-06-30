FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY . .
RUN go build -o app-service ./app

FROM alpine:latest

WORKDIR /root/
COPY --from=builder /app/app-service .

CMD ["./app-service"]