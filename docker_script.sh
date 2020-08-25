docker rm $(docker ps -a -q)
docker-compose -f docker-compose.yml up -d --build
