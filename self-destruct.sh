#!/bin/bash

# self-destruct script for easy-vpn virtual machine
# ==============================================================================
# this script will be uploaded and started in the background, on the virtual 
# machine that was created by running easy-vpn with the "up" command


# check arguments
if [[ $# -ne 4 ]]; then
    echo "usage: ./self-destruct.sh <API-Provider> <API-Key> <VM-Id> <Uptime>"
    exit 1
fi

PROVIDER=$1
API_KEY=$2
VM_ID=$3
UPTIME=$4

echo "Called ./self-destruct.sh ${PROVIDER} ${API_KEY} ${VM_ID} ${UPTIME}"

function destroyVM {
	if [[ $PROVIDER == "digitalocean" ]]; then
		echo "curl -X DELETE -i -H \"Authorization: Bearer ${API_KEY}\" \"https://api.digitalocean.com/v2/droplets/${VM_ID}\""
		curl -X DELETE -i -H "Authorization: Bearer ${API_KEY}" "https://api.digitalocean.com/v2/droplets/${VM_ID}"
		exit $?
	elif [[ $PROVIDER == "vultr" ]]; then
		echo "curl -X POST -i -d \"SUBID=${VM_ID}\" \"https://api.vultr.com/v1/server/destroy?api_key=${API_KEY}\""
		curl -X POST -i -d "SUBID=${VM_ID}" "https://api.vultr.com/v1/server/destroy?api_key=${API_KEY}"
		exit $?
	else
		echo "Unknown provider!"
		exit 5
	fi
}

STARTTIME=`date +%s`
while true; do
	CURRENTTIME=`date +%s`
	DEADLINE=$(( STARTTIME + UPTIME ))
	if [[ $CURRENTTIME -gt $DEADLINE ]]; then
		destroyVM
	fi

	sleep 5
done

exit 0
