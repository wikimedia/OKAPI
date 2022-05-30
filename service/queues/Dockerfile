###############
# Build stage #
###############
FROM golang:1.16 AS builder

# Switch to /app dir
WORKDIR /app

# Cache dependencies download
COPY go.mod go.sum ./
RUN go mod download

# Copy the project files
COPY . .

# Build app binary
RUN go build -o main queues/*.go

# Build migrations binary
RUN go build -o migrate migrations/*.go

#############
# Run stage #
#############
FROM alpine:latest

# Istall glibc compatibility
RUN apk add gcompat

# Switch to /root
WORKDIR /root

# Copy app binary from previous stage
COPY --from=builder /app/main ./

# Copy migrations binary from previous stage
COPY --from=builder /app/migrate ./

# Set the binary as the CMD of the container
CMD ./migrate migrate ; ./main
