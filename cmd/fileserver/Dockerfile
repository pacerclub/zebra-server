FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o fileserver ./cmd/fileserver

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/fileserver .
COPY --from=builder /app/.env .

EXPOSE 8081

CMD ["./fileserver"]
