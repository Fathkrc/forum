
# Stage 1: Build the Go application
FROM golang:1.26-bookworm AS builder

WORKDIR /src

# CGO is needed by SQLite drivers such as mattn/go-sqlite3.
RUN apt-get update && apt-get install -y --no-install-recommends \
    gcc \
    libc6-dev \
    && rm -rf /var/lib/apt/lists/*

# Copy dependency definitions first for better build caching.
COPY go.mod go.sum ./
RUN go mod download

# Copy the application source.
COPY . .

# Compile the application.
RUN CGO_ENABLED=1 go build -o /forum .

# Stage 2: Create a smaller runtime image.
FROM debian:bookworm-slim

WORKDIR /app

# Install certificates for HTTPS requests, if needed.
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Copy the compiled application.
COPY --from=builder /forum /app/forum

# Copy runtime assets (adjust these paths to your project).
COPY --from=builder /src/templates /app/templates
# Uncomment and adapt if your project has a static directory:
# COPY --from=builder /src/static /app/static

# Docker image metadata.
LABEL org.opencontainers.image.title="Literary Lions Forum" \
      org.opencontainers.image.description="A Go-based literary discussion forum" \
      org.opencontainers.image.version="1.0.0"

EXPOSE 8080

CMD ["/app/forum"]
