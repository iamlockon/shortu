#!bin/bash

set -euo pipefail

if [ "$#" -ne 1 ]; then
    echo "Usage: $0 <env>"
fi

ENV=$1

if [ -f conf/$ENV.env ]
then
  export $(cat conf/$ENV.env | sed 's/#.*//g' | xargs)
fi

if [ -f conf/$ENV.secrets.env ]
then
  export $(cat conf/$ENV.secrets.env | sed 's/#.*//g' | xargs)
fi

echo "env | grep PG:"
env | grep PG

docker-compose build
docker-compose up --force-recreate
