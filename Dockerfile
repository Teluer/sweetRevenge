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
RUN go build -gcflags \"all=-N -l\" -o sweetRevenge ./src/main

# Start a new stage
FROM rabbitmq:3.9

# Copy the built Go application from the previous stage
COPY --from=builder /app/sweetRevenge /usr/local/bin/sweetRevenge

COPY config.properties .
COPY rabbitmq.conf /etc/rabbitmq/

# Copy Tor and libraries
COPY TorLinux/data /usr/local/tor/data/
COPY TorLinux/tor/tor /usr/local/bin/
COPY TorLinux/torrc ./torrc
COPY TorLinux/tor/lib* /usr/local/lib/
RUN ldconfig /usr/local/lib/

# Copy the entrypoint script
COPY entrypoint.sh /usr/local/bin/entrypoint.sh

# Expose the port on which your Go application listens
EXPOSE 2345

# Set the entrypoint script as executable
RUN chmod +x /usr/local/bin/entrypoint.sh

# Set the command to run the entrypoint script
ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]