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
RUN go build -o sweetRevenge ./src/main

# Start a new stage
FROM rabbitmq:3.9

# Install necessary dependencies for Chrome
RUN apt-get update && apt-get install -y \
    ca-certificates \
    wget \
    gnupg

# Download and install Chrome
RUN wget -q -O - https://dl-ssl.google.com/linux/linux_signing_key.pub | apt-key add - && \
    echo "deb [arch=amd64] http://dl.google.com/linux/chrome/deb/ stable main" >> /etc/apt/sources.list.d/google-chrome.list && \
    apt-get update && \
    apt-get install -y google-chrome-stable

# Copy the built Go application from the previous stage
COPY --from=builder /app/sweetRevenge /usr/local/bin/sweetRevenge

COPY config.properties .
COPY user-agents .
COPY files/chromedriver /usr/local/bin/chromedriver
COPY rabbitmq.conf /etc/rabbitmq/

# Copy Tor and libraries
COPY files/TorLinux/data /usr/local/tor/data/
COPY files/TorLinux/tor/tor /usr/local/bin/
COPY files/TorLinux/torrc ./torrc
COPY files/TorLinux/tor/lib* /usr/local/lib/
RUN ldconfig /usr/local/lib/

# Copy the entrypoint script
COPY entrypoint.sh /usr/local/bin/entrypoint.sh

# Set the entrypoint script as executable
RUN chmod +x /usr/local/bin/entrypoint.sh

# Set the command to run the entrypoint script
ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]