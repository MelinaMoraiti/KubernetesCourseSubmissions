# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY . .

RUN go build -o log_output log_output.go

# Runtime stage
FROM scratch

WORKDIR /

COPY --from=builder /app/log_output /log_output

ENV PORT=3000
ENV PING_PONG_APP_HOST="http://localhost:8081"

EXPOSE 3000

ENTRYPOINT ["./log_output"]