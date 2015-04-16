#!/usr/bin/env bash

# declaring variables permanent that means they are stored forver
# you can't change them ever.
# you can see the readonly vairable
# readonly -p
# So once you start this script it will keep taking the same
# values as it was first ran
# Refer : http://bash.cyberciti.biz/guide/Readonly_command

# slapcat binary location
SLAPCAT='/usr/sbin/slapcat'

# s3cmd binary location
S3CMD='/usr/bin/s3cmd'

# s3 bucket where to backup
BUCKET='s3://<BUCKET_NAME>'

# get today's date
DATE=$(date +%F)

# ldif filename
filename="/tmp/ldapdb-${DATE}"

# slapd configuration file
SLAPDCONF='/etc/ldap/slapd.conf'

# options needs to pass to slapcat command
SLAPCATOPT="-v -l ${filename} -f ${SLAPDCONF}"

# getting script name for logging
PROGNAME=$(basename $0)


# logging messages to the system's standard location

log_message() {
    message=$1
    logger -i -s -t ${PROGNAME} ${message}
}

# check if s3cmd exists

if [[ -z ${S3CMD} ]];then
    printf "\ns3cmd not found.\nInstall s3cmd from http://s3tools.org/s3cmd\n"
    log_message "\ns3cmd not found.\nInstall s3cmd from http://s3tools.org/s3cmd\n"
    exit 1
fi

log_message "Taking dump from ldap server"

# build slapcat command and execute

${SLAPCAT} ${SLAPCATOPT}

if [[ $? -eq 1 ]];then
    printf "\n slapcat dump was not successful. Please run the command manually\
            and look for problem\n"
    log_message "${PROGNAME} was not successful"
    exit 1
else
    log_message "LDAP dump completed successful"
fi


# uploading the dump to s3
echo "${S3CMD} put ${filename} ${BUCKET}"

${S3CMD} put ${filename} ${BUCKET}

if [[ $? -eq 1 ]];then
    printf " S3 put command was not successful. Dump is not uploaded to S3\
        please  run command manually\n"
    log_message "${PROGNAME} S3 put command was not successful. Dump is not uploaded to S3"
    exit 1
else
    log_message "LDAP dump is uploaded to S3 successfully"
fi
