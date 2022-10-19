#!/bin/sh

docker stop carrot-backyard
docker rm carrot-backyard
docker rmi gongchen0618/carrot-backyard:0.0.1
docker rmi carrot-backyard

docker build -t carrot-backyard .
docker tag carrot-backyard gongchen0618/carrot-backyard:0.0.1
docker images|grep none|awk '{print $3}'|xargs docker rmi
