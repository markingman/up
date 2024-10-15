#!/bin/sh

# Upload config to VM in Google Cloud

source ../.env

for VAR in PROJECT HOST_ZONE NAME
	do [ -z "${!VAR}" ] && echo "\nERROR: $VAR not set\n" && exit 1
done

if [ -z "$1" ]
	then
	echo "\nERROR: set file name as \$1 \n" && exit 1
fi

FILE=$1
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

CMD="
sudo chown root:root ~/$BASE
sudo mv ~/$BASE /var/data/$BASE
"

gcloud compute ssh $NAME --project=$PROJECT --zone=$HOST_ZONE --tunnel-through-iap --command "$CMD"
