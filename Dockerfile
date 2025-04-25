# Start from the official Go image to build our application
FROM golang:1.24-alpine AS build

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod files first for better layer caching
COPY go.mod ./

# Copy the source code
COPY *.go ./

# Build the Go app
RUN go build -o server .

# Start a new stage from scratch for a smaller final image
FROM alpine:latest

WORKDIR /app

# Copy the binary from the build stage
COPY --from=build /app/server /app/

# Expose port 8080 to the outside
EXPOSE 8080

# Command to run the executable
CMD ["/app/server"]