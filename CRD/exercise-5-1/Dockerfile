# Build stage
FROM golang:1.26-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY *.go .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o site_controller .

# Final stage
FROM alpine:latest

RUN apk add --no-cache ca-certificates \
    && adduser -D -u 10001 app

COPY --from=builder /build/site_controller /usr/local/bin/site_controller

USER 10001
CMD ["site_controller"]