#!/bin/bash

# Deploy actions for deploying as a VM to Google Cloud

source .env

for VAR in NAME PROJECT HOST_ZONE HOST_REGION REPO SERVICE_ACCOUNT
	do [ -z "${!VAR}" ] && echo "\nERROR: $VAR not set (expected to be in .env file)\n" && exit 1
done

CMDS="deploy, create, reset, update, ssh"

if [ -z "$1" ]
    then echo "nERROR: No command provided. Available commands: $CMDS" && exit 1
fi

deploy() {
    echo "Build image for Google Cloud Artifact Registry $PROJECT / $REPO / $NAME..."

	TAG=$HOST_REGION-docker.pkg.dev/$PROJECT/$REPO/$NAME:latest

	docker build -f ../Dockerfile --platform linux/amd64 -t $TAG ../
	docker push $TAG
}

create() {
    echo "Create VM $PROJECT / $HOST_ZONE / $NAME..."

	gcloud compute instances create-with-container $NAME \
		--project=$PROJECT \
		--zone=$HOST_ZONE \
		--machine-type=f1-micro \
		--network-interface=network-tier=PREMIUM,subnet=default,no-address \
		--maintenance-policy=MIGRATE \
		--provisioning-model=STANDARD \
		--service-account=$SERVICE_ACCOUNT \
		--scopes=https://www.googleapis.com/auth/devstorage.read_only,https://www.googleapis.com/auth/logging.write,https://www.googleapis.com/auth/monitoring.write,https://www.googleapis.com/auth/service.management.readonly,https://www.googleapis.com/auth/servicecontrol,https://www.googleapis.com/auth/trace.append \
		--image=projects/cos-cloud/global/images/cos-stable-117-18613-0-76 \
		--boot-disk-size=10GB \
		--boot-disk-type=pd-balanced \
		--boot-disk-device-name=$NAME \
		--container-image=$HOST_REGION-docker.pkg.dev/$PROJECT/$REPO/$NAME:latest \
		--container-restart-policy=always \
		--container-command=/app/app \
		--container-mount-host-path=host-path=/var/data,mode=ro,mount-path=/app/data \
		--no-shielded-secure-boot \
		--shielded-vtpm \
		--shielded-integrity-monitoring \
		--labels=goog-ec-src=vm_add-gcloud,container-vm=cos-stable-117-18613-0-76 \
		--reservation-affinity=any \
		--no-address
}

reset() {
    echo "Resetting VM $PROJECT / $HOST_ZONE / $NAME..."

	gcloud compute instances reset $NAME --project=$PROJECT --zone=$HOST_ZONE
}

update() {
    echo "Updating data on VM $PROJECT / $HOST_ZONE / $NAME..."

	FILE="disk/conf.json"

	if [ -z "$FILE" ]
		then
		echo "\nERROR: could not find $FILE\n" && exit 1
	fi

	BASE=`basename $FILE`

	echo "Uploading..."

	gcloud compute scp \
		--project=$PROJECT \
		--zone=$HOST_ZONE \
		--recurse \
		--tunnel-through-iap \
		--compress \
		$FILE $NAME:~/

	echo "Copying to position..."

	CMD="sudo chown root:root ~/$BASE; sudo mv ~/$BASE /var/data/$BASE"

	gcloud compute ssh $NAME --project=$PROJECT --zone=$HOST_ZONE --tunnel-through-iap --command "$CMD"
}

ssh() {
    echo "SSH to VM $PROJECT / $HOST_ZONE / $NAME..."

	gcloud compute ssh $NAME --project=$PROJECT --zone=$HOST_ZONE
}

case "$1" in
    deploy)
        deploy
        ;;
    create)
        create
        ;;
    reset)
        reset
        ;;
    update)
        update
        ;;
    ssh)
        ssh
        ;;
    *)
        echo "nERROR: Invalid command. Available commands: $CMDS" && exit 1
        ;;
esac
