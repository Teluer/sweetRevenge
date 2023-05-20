# Use the official Golang image as the base
FROM golang:1.19 as builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files
COPY go.mod go.sum ./

# Download the Go module dependencies
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go application
RUN go build -o /usr/local/bin/sweetRevenge ./src/main

# Copy the Tor binary
COPY TorLinux /usr/local/bin/
# Copy the required library for Tor
COPY TorLinux/tor/lib* /usr/local/lib/
# Update the library cache
RUN ldconfig /usr/local/lib/

# Copy the entrypoint script
COPY entrypoint.sh /usr/local/bin/entrypoint.sh

# Expose the port on which your Go application listens
#EXPOSE 8080

# Set the entrypoint script as executable
RUN chmod +x /usr/local/bin/entrypoint.sh

# Set the command to run the entrypoint script
ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]