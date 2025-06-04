FROM golang:1.23.4-alpine3.21 AS builder

COPY . /auth-service/
WORKDIR /auth-service/

RUN go mod download
RUN go build -o ./bin/auth-server cmd/server/main.go

FROM alpine:3.21
WORKDIR /root/
COPY --from=builder /auth-service/bin/auth-server .
COPY .env .

CMD ["./auth-server"]
