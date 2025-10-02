FROM golang:1.21-alpine

WORKDIR /app

# VULNERABILITY: Running as root
RUN apk add --no-cache curl

# Copy pre-built binary and dependencies
COPY cmd/api/main ./
COPY go.mod go.sum* ./

# Create directories for file operations
RUN mkdir -p /app/uploads /app/documents /backups

EXPOSE 8082

# VULNERABILITY: Running as root allows full system access
CMD ["./main"]
