# Use the official Golang image to create a build artifact.
FROM golang:1.24 AS builder

WORKDIR /app

# Retrieve application dependencies.
COPY go.mod go.sum ./
RUN go mod download

# Copy local code to the container image.
COPY . .

# Build the binary.
RUN CGO_ENABLED=0 GOOS=linux go build -v -o server ./cmd/server

# Use a lean image for the runtime.
FROM alpine:3.19
RUN apk add --no-cache ca-certificates

WORKDIR /app

# Copy the binary to the production image from the builder stage.
COPY --from=builder /app/server /app/server

# Expose port 8080 to the outside world.
EXPOSE 8080

# Run the web service on container startup.
CMD ["/app/server"]
