#!/bin/bash

# Start the Tor service
/usr/local/bin/tor -f torrc &> tor.log &

#Start rabbitmq
rabbitmq-server &> rabbit.log &

# Run your Go application
dlv --listen=:2345 --headless=true --api-version=2 --accept-multiclient exec ./sweetRevenge