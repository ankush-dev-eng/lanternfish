# Stage 1: Build the binary
FROM golang:1.22-alpine AS builder

WORKDIR /build

# Install ca-certificates and tzdata
RUN apk --no-cache add ca-certificates tzdata

# Copy dependency manifests
COPY go.mod ./
RUN go mod download

# Copy source code
COPY . .

# Build statically linked standalone binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o lanternfish ./cmd/lanternfish

# Stage 2: Minimal runtime image
FROM alpine:3.20

WORKDIR /app

RUN apk --no-cache add ca-certificates tzdata

# Copy binary from builder
COPY --from=builder /build/lanternfish /app/lanternfish

# Default environment
ENV PORT=8080
EXPOSE 8080

# Run as non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser

ENTRYPOINT ["/app/lanternfish"]
