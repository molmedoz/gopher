# Gopher - Go Version Manager
# Minimal Docker image for running Gopher in containerized environments

FROM golang:1.24-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make

# Set working directory
WORKDIR /build

# Copy source code
COPY . .

# Build gopher binary
RUN go build -ldflags="-s -w" -o gopher ./cmd/gopher

# Create final minimal image
FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache \
    ca-certificates \
    bash \
    curl \
    git

# Copy gopher binary from builder
COPY --from=builder /build/gopher /usr/local/bin/gopher

# Create non-root user for security
RUN adduser -D -s /bin/bash -h /home/gopher gopher

# Switch to non-root user
USER gopher
WORKDIR /home/gopher

# Set up Gopher environment
ENV HOME="/home/gopher"
ENV PATH="/home/gopher/.local/bin:$PATH"

# Create Gopher directories with proper permissions
RUN mkdir -p .gopher/versions .gopher/downloads .local/bin

# Verify gopher works
RUN gopher version

# Default entrypoint
ENTRYPOINT ["gopher"]
CMD ["help"]

