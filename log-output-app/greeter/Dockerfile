# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY . .

RUN go build -o greeter greeter.go

# Runtime stage
FROM scratch

WORKDIR /

COPY --from=builder /app/greeter /greeter

ENV PORT=5000

EXPOSE 5000

ENTRYPOINT ["./greeter"]