#!/bin/bash

# Start the Tor service
/usr/local/bin/tor -f torrc &> tor.log &

# Run your Go application
sweetRevenge