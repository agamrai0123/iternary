# Multi-stage build for minimal image
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o itinerary-backend .

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates postgresql-client

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/itinerary-backend .
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/static ./static
COPY --from=builder /app/config ./config

# Expose port
EXPOSE 8080

# Run the application
CMD ["./itinerary-backend"]
