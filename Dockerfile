# Stage 1: Build Stage
FROM ubuntu:20.04 AS builder

# Set noninteractive mode for apt-get
ENV DEBIAN_FRONTEND=noninteractive

# Install required packages: wget, tar, git, ca-certificates
RUN apt-get update && apt-get install -y \
    wget \
    tar \
    git \
    ca-certificates \
 && rm -rf /var/lib/apt/lists/*

# Install Go (adjust version as needed)
ENV GO_VERSION=1.21.1
RUN wget https://golang.org/dl/go${GO_VERSION}.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz && \
    rm go${GO_VERSION}.linux-amd64.tar.gz

# Set Go environment variables
ENV PATH="/usr/local/go/bin:${PATH}"
WORKDIR /app

# Copy go.mod and go.sum first to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY . .

# Create .env file from .env.example
RUN cp .env.example .env

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -o backend ./cmd/main.go

# Stage 2: Final Image
FROM ubuntu:20.04

# Install Redis client
RUN apt-get update && apt-get install -y \
    ca-certificates \
    redis-tools \
 && rm -rf /var/lib/apt/lists/*

# Copy the built binary and env file from the builder stage
COPY --from=builder /app/backend /backend
COPY --from=builder /app/.env /.env

# Expose port 8080 (or whichever port your app uses)
EXPOSE 8080

# Command to run the binary
ENTRYPOINT ["/backend"]
