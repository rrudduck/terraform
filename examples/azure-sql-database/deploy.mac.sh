#!/bin/bash

set -o errexit -o nounset

if docker -v; then

  # generate a unique string for CI deployment
  export KEY=$(cat /dev/urandom | env LC_CTYPE=C tr -cd 'a-z' | head -c 12)
  export PASSWORD=$a@abcd9753w0w@h@12
  # =$KEY$(cat /dev/urandom | env LC_CTYPE=C tr -cd 'A-Z' | head -c 2)$(cat /dev/urandom | env LC_CTYPE=C tr -cd '0-9' | head -c 2)

  /bin/sh ./deploy.ci.sh

else
  echo "Docker is used to run terraform commands, please install before run:  https://docs.docker.com/docker-for-mac/install/"
fi