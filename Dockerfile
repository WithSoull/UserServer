FROM golang:1.24.1-alpine3.21 AS builder

COPY . /user-server/
WORKDIR /user-server/

RUN go mod download
RUN go build -o ./bin/user-server cmd/server/main.go

FROM alpine:3.21
WORKDIR /root/
COPY --from=builder /user-server/bin/user-server .
COPY .env .
COPY service.pem .
COPY service.key .
COPY ca.cert .


CMD ["./user-server"]
