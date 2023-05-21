docker build -t sweet-revenge:latest .
docker run --network="host" -p 2345:2345 --restart unless-stopped sweet-revenge