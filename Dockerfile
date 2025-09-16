# Stage 1: Build the Go binary
FROM golang:1.24.7 AS builder

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum to cache dependencies
COPY go.mod go.sum ./

# Download dependencies (if any)
RUN go mod download

# Copy the source code
COPY cmd/ ./cmd/

# Build the binary with optimizations
RUN CGO_ENABLED=0 GOOS=linux go build -o /testflow ./cmd

# Stage 2: Create a minimal runtime image
FROM scratch

# Copy the binary from the builder stage
COPY --from=builder /testflow /testflow

# Set the entrypoint to run the CLI
ENTRYPOINT ["/testflow"]