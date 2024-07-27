# Start from the latest official Go image
FROM golang:1.22 as builder

# Set the working directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Start a new stage from scratch
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the pre-built binary file from the previous stage
COPY --from=builder /app/main .

# Copy any additional necessary files (like static assets)
COPY --from=builder /app/migrations ./migrations
# COPY --from=builder /app/pb_public ./pb_public

# Command to run the executable
CMD ["./main", "serve", "--http=0.0.0.0:8080"]
