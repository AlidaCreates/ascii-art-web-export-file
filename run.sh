docker image build -f Dockerfile -t ascii-art-web-docker .

docker images

docker container run -p 8080:8080 --detach --name dockerize ascii-art-web-docker

docker ps -a

docker exec -it dockerize /bin/sh -c "ls -l"