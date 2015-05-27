#!/usr/bin/env bash

# s3 binary location
readonly S3CMD='/usr/bin/s3cmd'

# optionns needs to pass to the s3cmd
readonly S3CMD_OPTIONS='sync --recursive --skip-existing'

# curl binary location
readonly CURL='/usr/bin/curl'

# get the instnace id
readonly INSTANCE_ID=$(${CURL} http://169.254.169.254/latest/meta-data/instance-id)

# s3 bucket where to backup
readonly BUCKET="s3://foo/application_logs/${INSTANCE_ID}/"

# source folder to backup
readonly BACKUP_PATH='/mnt/ephemeral/logs/'


# getting script name for logging
readonly PROGNAME=$(basename $0)

# how many old days files need to backup
readonly AGE=3


# logging messages to the system's standard location

log_message() {
    message=$1
    logger -i -s -t ${PROGNAME} ${message}
}


deletefiles() {
# find all the files for upload

TODELETE=$(find ${BACKUP_PATH} -type f -mtime +${AGE})

log_message "Deleting files older than ${AGE}"

if test -z "${TODELETE}";then
        log_message "No files to delete older than ${AGE}"
else
        for file in "${TODELETE[@]}";do
                rm $file
        done
fi
}


# check if s3cmd exists

if [[ -z ${S3CMD} ]];then
    printf "\ns3cmd not found.\nInstall s3cmd from http://s3tools.org/s3cmd\n"
    exit 1
fi

# check if source directory exists

if [[ ! -d ${BACKUP_PATH} ]];then
    printf "\n source backup path ${BACKUP_PATH} doesnt exists\n"
    log_message "${BACKUP_PATH} does not exists"
    exit 1
fi



#uploading all the files to s3

log_message "Syncing ${BACKUP_PATH} to ${BUCKET}"

${S3CMD} ${S3CMD_OPTIONS} ${BACKUP_PATH} ${BUCKET}

# if sync was not successful dont delete

if [[ $? -eq 1 ]];then
    log_message "Sync was not successful"
    exit 1
else
    deletefiles
fi
