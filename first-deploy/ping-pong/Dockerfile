# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /build

# Copy go mod files
COPY go.mod ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ping-pong pingpong.go

# Final stage
FROM scratch

WORKDIR /

# Copy binary from builder
COPY --from=builder /build/ping-pong .
ENV PORT=9000
# Expose port
EXPOSE 9000

# Run the application
CMD ["./ping-pong"]