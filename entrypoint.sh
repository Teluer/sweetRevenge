#!/bin/bash

# Start the Tor service
/usr/local/bin/tor/tor -f /usr/local/bin/torrc &> tor.log &

# Run your Go application
sweetRevenge