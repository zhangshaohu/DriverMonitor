#!/bin/bash

set -e

# Deployment helper script which should run on the server we're deploying to.

sudo apt-get update

# install docker
sudo apt-get install -y \
     apt-transport-https \
     ca-certificates \
     curl \
     gnupg2 \
     software-properties-common

curl -fsSL https://download.docker.com/linux/debian/gpg | sudo apt-key add -

sudo add-apt-repository \
     "deb [arch=amd64] https://download.docker.com/linux/debian \
   $(lsb_release -cs) \
   stable"

sudo apt-get update

sudo apt-get install -y docker-ce

sudo gpasswd -a $(whoami) docker

# TODO this will fail because we need to be in the docker group to run --
# will work on second login though

# install docker-compose
sudo curl -L https://github.com/docker/compose/releases/download/1.21.0/docker-compose-$(uname -s)-$(uname -m) -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# pull api image

gcloud beta auth configure-docker -q
docker pull us.gcr.io/delta-geode-199714/csc591-api
docker pull us.gcr.io/delta-geode-199714/csc591-db-watcher

source .envrc

# start compose services
docker-compose down
docker volume create mosquitto_data || true
docker volume create mosquitto_log || true
docker volume rm influx_data || true
docker volume create influx_data || true
docker-compose up -d

# load data from dump
docker exec -t csc591_influxdb_1 bash -c "curl -XPOST -d @/influxdata.txt 'http://localhost:8086/write?db=csc591' && influx -import -path /influxdata.txt"
