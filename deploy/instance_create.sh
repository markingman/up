#!/bin/sh

# Create a VM in Google Cloud
# Edit the `gcloud compute instances create` options below as required.

source ../.env

for VAR in NAME PROJECT HOST_ZONE HOST_REGION REPO SERVICE_ACCOUNT
	do [ -z "${!VAR}" ] && echo "\nERROR: $VAR not set\n" && exit 1
done

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
