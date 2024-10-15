#!/bin/sh

source ../.env

for VAR in NAME PROJECT HOST_ZONE
	do [ -z "${!VAR}" ] && echo "\nERROR: $VAR not set\n" && exit 1
done

gcloud compute ssh $NAME --project=$PROJECT --zone=$HOST_ZONE
