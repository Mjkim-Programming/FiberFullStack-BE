FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o fiber-app main.go

FROM ubuntu:22.04

WORKDIR /root/

COPY --from=builder /app/fiber-app .

EXPOSE 4000

CMD ["./fiber-app"]