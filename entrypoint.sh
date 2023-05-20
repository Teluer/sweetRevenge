#!/bin/bash

# Start the Tor service
tor &

# Start the MySQL database
service mysql start

# Run your Go application
myapp