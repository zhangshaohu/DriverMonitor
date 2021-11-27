#!/bin/bash

set -e

if [[ $# -ne 1 ]]; then
    echo "Usage: $0 <server SSH connection string>"
    echo " ex. $0 user@server.com"
    exit 1
fi

server_ssh_string="$1"

# Create optimized production build
docker-compose run --rm webpack yarn build

# Tag and push API image to the registry
gcloud_tag_base="us.gcr.io/delta-geode-199714"

api_image='csc591-api'
api_gcloud_tag="$gcloud_tag_base/$api_image"
docker tag "$api_image" "$api_gcloud_tag"
docker push "$api_gcloud_tag"

watcher_image='csc591-db-watcher'
watcher_gcloud_tag="$gcloud_tag_base/$watcher_image"
docker tag "$watcher_image" "$watcher_gcloud_tag"
docker push "$watcher_gcloud_tag"

# Export sample data to load
docker-compose up -d influxdb
docker-compose exec influxdb influx_inspect export \
               -datadir "/var/lib/influxdb/data" \
               -waldir "/var/lib/influxdb/wal" \
               -out "/influxdata.txt" \
               -database __test_csc591
docker cp csc591iotprojectwebapp_influxdb_1:/influxdata.txt deploy/influxdata.txt
sed -i "s/__test_csc591/csc591/g" deploy/influxdata.txt

# Copy up deploy script, compose script, nginx config, static files
project_dir="/opt/csc591"
function rsync_to_server() {
    rsync -rvz "$1" "$server_ssh_string:$project_dir"
}

ssh "$server_ssh_string" "sudo mkdir -p $project_dir && sudo chown \$(whoami) $project_dir"
ssh "$server_ssh_string" "sudo apt-get update && sudo apt-get install rsync"
rsync_to_server deploy/deploy_on_server.sh
rsync_to_server deploy/docker-compose.yml
rsync_to_server deploy/csc591.conf
rsync_to_server deploy/mosquitto.conf
rsync_to_server deploy/mosquitto_passwd
rsync_to_server deploy/influxdata.txt
rsync_to_server deploy/.envrc
rsync_to_server frontend/build

# Run the setup on server script
ssh "$server_ssh_string" "cd $project_dir && ./deploy_on_server.sh"
