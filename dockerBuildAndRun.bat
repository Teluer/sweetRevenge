docker build -t sweet-revenge:latest .
docker run -p 8008:8008 --restart unless-stopped sweet-revenge