#!/bin/sh

# Build image for Google Cloud Artifact Registry

source ../.env

for VAR in HOST_REGION PROJECT REPO NAME
	do [ -z "${!VAR}" ] && echo "\nERROR: $VAR not set\n" && exit 1
done

TAG=$HOST_REGION-docker.pkg.dev/$PROJECT/$REPO/$NAME:latest

docker build -f ../Dockerfile --platform linux/amd64 -t $TAG ../
docker push $TAG
