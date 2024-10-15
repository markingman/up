#!/bin/bash

# Reset a VM in Google Cloud

source ../.env

for VAR in PROJECT HOST_ZONE NAME
	do [ -z "${!VAR}" ] && echo "\nERROR: $VAR not set\n" && exit 1
done

gcloud compute instances reset $NAME --project=$PROJECT --zone=$HOST_ZONE
