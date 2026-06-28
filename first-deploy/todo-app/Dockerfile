# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o todoapp main.go

# Final stage
FROM scratch

WORKDIR /

# Copy binary from builder
COPY --from=builder /build/todoapp .

# Expose port
EXPOSE 8000

# Run the application
CMD ["./todoapp"]