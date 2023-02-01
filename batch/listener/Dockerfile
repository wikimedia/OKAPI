###############
# Build stage #
###############
FROM golang:1.18 AS builder

# Switch to /app dir
WORKDIR /app

# Cache dependencies download
COPY go.mod go.sum ./
RUN go mod download

# Copy the project files
COPY . .

# Build the binary
RUN go build -o main listener/*.go

#############
# Run stage #
#############
FROM alpine:latest

# Istall glibc compatibility
RUN apk add gcompat

# Switch to /root
WORKDIR /root

# Copy binary from previous stage
COPY --from=builder /app/main ./

# Set the binary as the CMD of the container
CMD ["./main"]
