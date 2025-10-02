FROM golang:1.21-alpine

WORKDIR /app

# VULNERABILITY: Running as root
# VULNERABILITY: No security updates

# Copy pre-built binary and dependencies
COPY cmd/api/main ./
COPY go.mod go.sum* ./

EXPOSE 8081

# VULNERABILITY: No non-root user
CMD ["./main"]
