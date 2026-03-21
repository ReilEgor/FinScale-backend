FROM golang:1.26-alpine AS builder
WORKDIR /app

COPY services/TransactionService/go.mod ./services/TransactionService/

WORKDIR /app/services/TransactionService
RUN go mod download

WORKDIR /app
COPY . .

RUN go build -o /transaction-service-app ./services/TransactionService/cmd

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/
COPY --from=builder /transaction-service-app .
CMD ["./transaction-service-app"]