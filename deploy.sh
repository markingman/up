#!/bin/bash

# Deploy actions for deploying as an image to Google Cloud Artifact Registry

source .env

for VAR in PROJECT HOST_REGION REPO TAG
	do [ -z "${!VAR}" ] && echo "\nERROR: $VAR not set (expected to be in .env file)\n" && exit 1
done

TAG=$HOST_REGION-docker.pkg.dev/$PROJECT/$REPO/$TAG

build() {
    echo "Build image for Google Cloud Artifact Registry $PROJECT / $REPO / $TAG..."

	docker build --platform linux/amd64 -t $TAG .
}

push() {
    echo "Push image to Google Cloud Artifact Registry $PROJECT / $REPO / $TAG..."

	docker push $TAG
}

deploy() {
    echo "Build and push image to Google Cloud Artifact Registry $PROJECT / $REPO / $TAG..."

	build
	push
}

case "$1" in
    "--build")
        build
        ;;
    "--push")
        push
        ;;
    "--deploy")
        deploy
        ;;
    *)
		CMDS="
--build
--push
--deploy
"

		if [ -z "$1" ]; then
				printf "\nERROR: No command provided\n"
			else
				printf "\nERROR: Invalid command\n":
		fi

		printf "\nAvailable commands: $CMDS\n" && exit 1
        ;;
esac
