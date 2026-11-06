# Build stage
FROM golang:1.26-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o broadcaster ./cmd

# Final stage
FROM alpine:latest

RUN adduser -D -u 10001 app
COPY --from=builder /build/broadcaster /broadcaster

ENV PORT=8080
EXPOSE 8080
USER app

CMD ["/broadcaster"]
