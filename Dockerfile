# Gunakan image Golang resmi sebagai builder
FROM golang:1.24-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go.mod dan go.sum lalu download library
COPY go.mod go.sum ./
RUN go mod download

# Copy seluruh source code
COPY . .

# Build aplikasi ke dalam binary bernama "main"
RUN go build -o main ./cmd/main.go

# Gunakan image alpine yang sangat ringan untuk menjalankan binary
FROM alpine:latest
WORKDIR /app

# Tambahkan ini untuk koneksi database luar (SSL/TLS)
RUN apk --no-cache add ca-certificates

# Copy binary dari stage builder
COPY --from=builder /app/main .
# Copy file .env (opsional, tapi biasanya disetel di docker-compose)
COPY .env .

# Jalankan aplikasi
CMD ["./main"]