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
RUN go build -o main *.go

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

# Copy model from previous stage
COPY --from=builder /app/model.conf ./

# Copy policy from previous stage 
COPY --from=builder /app/policy.csv ./

# Set the binary as the CMD of the container
CMD ["./main"]
