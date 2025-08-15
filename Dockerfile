# ---- Build Stage ----
# Use the official Golang image to create a build environment.
FROM golang:1.24-alpine AS builder

ENV CGO_ENABLED=0
ENV GOOS=linux

# Set the working directory inside the container.
WORKDIR /app

# Copy go.mod and go.sum first to leverage Docker's layer caching.
# This step will only be re-run if these files change.
COPY go.mod go.sum ./

# Tidy can fix inconsistencies between go.mod and go.sum.
RUN go mod tidy
# Verify the integrity of the dependencies.
RUN go mod verify
# Download the dependencies.
RUN go mod download

# Copy the rest of the application source code.
COPY . .

# ENV CGO_ENABLED=0
# ENV GOOS=linux

# # Set the working directory inside the container.
# WORKDIR /app

# # Copy go.mod
# COPY . .
# RUN go mod tidy
# RUN go mod install

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/server ./cmd/server

# ---- Production Stage ----
# Use a minimal base image. 'scratch' is the smallest possible image,
# but 'alpine' is a good choice if you need a shell for debugging.
FROM alpine:latest

# It's a good practice to run as a non-root user.
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Copy the compiled binary from the builder stage.
COPY --from=builder /app/server /app/server
COPY .env .

# Copy any other necessary assets, like configuration files or templates if needed.
# COPY --from=builder /app/config.json /app/config.json

# Set the user to the non-root user.
USER appuser

# Expose the port the application will run on.
EXPOSE 4200

# The command to run the application.
ENTRYPOINT ["/app/server"]
