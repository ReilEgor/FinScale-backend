FROM golang:1.26-alpine AS builder
WORKDIR /app

COPY services/UserService/go.mod ./services/UserService/

WORKDIR /app/services/UserService
RUN go mod download

WORKDIR /app
COPY . .

RUN go build -o /user-service-app ./services/UserService/cmd

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/
COPY --from=builder /user-service-app .
CMD ["./user-service-app"]