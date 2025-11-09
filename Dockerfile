########################################
# Stage 1: Build the Go binary
########################################
FROM golang:1.25.0-alpine AS builder

WORKDIR /app

# Build with a static binary for a minimal runtime image
ENV CGO_ENABLED=0 GOOS=linux

# Cache modules
COPY go.mod ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the API binary
RUN go build -trimpath -buildvcs=false -ldflags="-s -w" -o /app/studentapi ./cmd/studentapi

########################################
# Stage 2: Run the app
########################################
FROM alpine:3.20

WORKDIR /app

# Install CA certificates (for HTTPS requests, etc.) and add non-root user with home
RUN apk add --no-cache ca-certificates \
	&& adduser -D -u 10001 -h /home/appuser appuser \
	&& mkdir -p /home/appuser \
	&& chown -R appuser:appuser /home/appuser

# Copy the compiled binary from builder stage
COPY --from=builder /app/studentapi /app/studentapi

# Optionally copy config (non-fatal if not used in container)
# Uncomment if you want to ship local config files inside the image
# COPY --from=builder /app/config /app/config

# Expose API port
EXPOSE 8080

# Run in release mode; set HOME to avoid cache path issues
ENV GIN_MODE=release \
	HOME=/home/appuser \
	XDG_CACHE_HOME=/home/appuser/.cache
USER appuser

ENTRYPOINT ["/app/studentapi"]
