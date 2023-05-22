docker build -t sweet-revenge:latest .
docker run --network="host" --restart unless-stopped sweet-revenge