#!/usr/bin/env bash
set -e

NETWORK=wp_net
DB_CONTAINER=wp_db
WP_CONTAINER=wp_app
DB_VOLUME=wp_db_data
WP_VOLUME=wp_wp_data

DB_NAME=wordpress
DB_USER=wp
DB_PASS=wp
DB_ROOT_PASS=root

# network
if ! docker network inspect $NETWORK >/dev/null 2>&1; then
  docker network create $NETWORK
fi

# volume
docker volume inspect $DB_VOLUME >/dev/null 2>&1 || docker volume create $DB_VOLUME
docker volume inspect $WP_VOLUME >/dev/null 2>&1 || docker volume create $WP_VOLUME

# db
if ! docker inspect $DB_CONTAINER >/dev/null 2>&1; then
  docker run -d \
    --name $DB_CONTAINER \
    --network $NETWORK \
    -v $DB_VOLUME:/var/lib/mysql \
    -e MYSQL_DATABASE=$DB_NAME \
    -e MYSQL_USER=$DB_USER \
    -e MYSQL_PASSWORD=$DB_PASS \
    -e MYSQL_ROOT_PASSWORD=$DB_ROOT_PASS \
    mysql:8
fi

# wordpress
if ! docker inspect $WP_CONTAINER >/dev/null 2>&1; then
  docker run -d \
    --name $WP_CONTAINER \
    --network $NETWORK \
    -p 8090:80 \
    -v $WP_VOLUME:/var/www/html \
    -e WORDPRESS_DB_HOST=$DB_CONTAINER:3306 \
    -e WORDPRESS_DB_NAME=$DB_NAME \
    -e WORDPRESS_DB_USER=$DB_USER \
    -e WORDPRESS_DB_PASSWORD=$DB_PASS \
    wordpress:latest
else
  docker start $WP_CONTAINER
  docker start $DB_CONTAINER
fi
