FROM golang:1.26-alpine AS builder
WORKDIR /app

COPY services/AuthValidatorService/go.mod ./services/AuthValidatorService/

WORKDIR /app/services/AuthValidatorService
RUN go mod download

WORKDIR /app
COPY . .

RUN go build -o /auth-validator-service-app ./services/AuthValidatorService/cmd

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/
COPY --from=builder /auth-validator-service-app .
CMD ["./auth-validator-service-app"]