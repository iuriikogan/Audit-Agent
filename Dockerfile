# Use a 2-stage build for production efficiency
# This Dockerfile expects web/out to be present in the build context.

# --- Stage 1: Backend Builder ---
FROM golang:1.25 AS backend-builder
WORKDIR /app
ENV GOTOOLCHAIN=auto

# Copy Go dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build arguments
ARG TARGET=server

# If building the server, ensure the frontend assets are in the correct place for go:embed
RUN if [ "$TARGET" = "server" ]; then \
      mkdir -p cmd/server/out && \
      if [ -d "web/out" ]; then \
        cp -r web/out/* cmd/server/out/; \
      else \
        echo "Warning: web/out not found. Server will build without embedded assets or fallback to empty." >&2; \
      fi \
    fi

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/app ./cmd/${TARGET}/main.go

# --- Stage 2: Final Runtime Image ---
FROM alpine:latest
WORKDIR /app

# Install ca-certificates for external API calls
RUN apk --no-cache add ca-certificates

# Copy the binary from the backend builder
COPY --from=backend-builder /app/bin/app ./app

# Expose port
EXPOSE 8080

# Set entrypoint
CMD ["./app"]
