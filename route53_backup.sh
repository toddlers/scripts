#!/bin/bash

# This requires cli53 command utility
# please install using below command
# sudo pip install cli53

#Declare backup path and master zone files
BACKUP_PATH="$(date +%F)"
ZONES_FILE="all-zones.txt"
DNS_FILE="all-dns.txt"

#create date-stamped backup directory and enter it
mkdir -p "$BACKUP_PATH"
cd "$BACKUP_PATH"

# Create a list of all hosted zones
cli53 list > "$ZONES_FILE" 2>&1

#create a list of domain names only
sed '/Name:/!d' "$ZONES_FILE"|cut -d: -f2|sed 's/^.//'|sed 's/.$//' > "$DNS_FILE"

#creat backup files for each domain
while read -r line; do
    cli53 export --full "$line" > "$line.txt"
done < "$DNS_FILE"

tar cvfz "$BACKUP_PATH.tgz" "$BACKUP_PATH"


exit 0

