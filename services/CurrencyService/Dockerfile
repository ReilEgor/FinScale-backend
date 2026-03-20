FROM golang:1.26-alpine AS builder
WORKDIR /app

COPY services/CurrencyService/go.mod ./services/CurrencyService/

WORKDIR /app/services/CurrencyService
RUN go mod download

WORKDIR /app
COPY . .

RUN go build -o /currency-service-app ./services/CurrencyService/cmd

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/
COPY --from=builder /currency-service-app .
CMD ["./currency-service-app"]