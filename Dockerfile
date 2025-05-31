FROM golang:1.24 AS builder

WORKDIR /app
COPY . .

RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./cmd/api

FROM debian:bullseye-slim
WORKDIR /app
COPY --from=builder /server ./

EXPOSE 8080
CMD ["./server"]