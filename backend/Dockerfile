FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main ./cmd/app

FROM alpine:3.19

WORKDIR /root/

COPY --from=builder /app/main .

RUN mkdir -p /root/uploads

EXPOSE 8080

CMD ["./main"]
