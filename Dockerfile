# ---------- Build Stage ----------
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install git (needed for private modules sometimes)
RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o log-collector

# ---------- Runtime Stage ----------
FROM alpine:latest

WORKDIR /app

# Add CA certificates for HTTPS connections
RUN apk --no-cache add ca-certificates

COPY --from=builder /app/log-collector .

EXPOSE 50051

ENV OTLP_ENDPOINT=localhost:4318

CMD ["./log-collector"]
