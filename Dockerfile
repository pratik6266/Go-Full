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

# Install CA certificates (for HTTPS requests, etc.) and add non-root user
RUN apk add --no-cache ca-certificates \
	&& adduser -D -H -u 10001 appuser

# Copy the compiled binary from builder stage
COPY --from=builder /app/studentapi /app/studentapi

# Expose API port
EXPOSE 8080

# Run in release mode and drop privileges
ENV GIN_MODE=release
USER appuser

ENTRYPOINT ["/app/studentapi"]
