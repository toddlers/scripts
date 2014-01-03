#!/usr/bin/env bash

# TODO
# Have a config file for taking these params 
# Have a single user to do all the task
# Add the feature for taking the VCS from the commandline 
# Add feature to push the crontab files to S3
# Add checks for svn status check as well.

TOPDIR=<TOPDIR>
SLDIR=$TOPDIR/Softlayer
AWSDIR=$TOPDIR/AWS
RUNAS=<USERNAME>
KEY=<KEY_FILE_FOR_SSH>
EMAIL=crons@example.com
SVNLOCATION='http://svn.corp.example.com/svn/prod-cronjobs'
VCS=svn

function show_help() {
cat << _USAGE
usage: $0 [-h] [-v VCS] [-t TOPDIR] [-u USER]
    -h          Show help
    -t TOPDIR   Set TOPDIR as top level dir, default $TOPDIR
    -u USER     Run as user USER, default $RUNAS
    -v VCS      Provide the VCS name, default $VCS
_USAGE
exit 0
}

while getopts "ht:u:v:" opt; do
    case "$opt" in
        h) show_help
            ;;
        t) TOPDIR=$OPTARG
            ;;
        u) RUNAS=$OPTARG
            ;;
        v) VCS=$OPTARG
            ;;
    esac
done

function _fmt() {
local color_ok="\x1b[32m"
local color_bad="\x1b[31m"
}

if [ $(id -u -n) != "${RUNAS}" ]; then
    echo "Please run as user ${RUNAS}"
    exit 1
fi

echo "Creating Directories for Separate crontabs from SERVERSIP and AWS servers"
mkdir -p $SLDIR $AWSDIR

# IP of the AWS servers

AWS_SERVERS=(IP1 IP2)


SL_SERVERS=(IP1 IP2)


    echo "Copying Crons from the AWS servers $AWSDIR"
for aserver in "${AWS_SERVERS[@]}"; do
    ssh -A -i $KEY root@$aserver 'crontab -l' > $AWSDIR/$aserver.$(date +%Y%m%d)
done


    echo "Copying crons from Softlayer servers in $SLDIR"
for server in "${SL_SERVERS[@]}"; do
    ssh -p2223 -t $RUNAS@$server 'sudo cat /var/spool/cron/root' > $SLDIR/$server.$(date +%Y%m%d) 2> /dev/null
done

cd $TOPDIR


# Using svn as vcs provider for taking the backup

function svn_backup() {

    changed_content=`svn st`

    echo "Adding all the contents in svn"

    svn add $SLDIR $AWSDIR

    if [[ $? != 0 ]]; then
        echo "snv add command failed while taking crontabs backup" | mail -s "[CRONTABS_BACKUP] svn add failed" $EMAIL
    exit 1
    fi

    svn ci --username rightster --password ZLwHXMsp  -m "[CRONTABS_BACKUP ] Commiting the crons from production servers $(date +%Y%m%d)"

    if [[ $? != 0 ]]; then
        echo "svn commit failed while taking the crontabs backup"| mail -s "[CRONTABS_BACKUP] svn commit failed" $EMAIL
    exit 1
    fi

    if [ -z "${changed_content}" ]; then
        echo "There is no change in crontabs. Please find the crontabs svn location ${SVNLOCATION}"| mail -s "[CRONTABS_BACKUP] Committedthe crons from producion servers" $EMAI
    else
        echo "Please find the crontabs svn location ${SVNLOCATION} ${changed_content}"| mail -s "[CRONTABS_BACKUP] Committedi the crons from producion servers" $EMAIL
    fi
}

# Using git as vcs provider for taking the backup

function git_backup {

    git add *

    if [[ $? != 0 ]] then
        echo "while doing crontabs backup , git add command failed"| mail -s "[CRONTABS_BACKUP] git add failed" $EMAIL
        exit 1
    fi

    git commit -a -m "[CRONTABS_BACKUP ] Commiting the crons from production servers $(date +%Y%m%d)"

    if [[ $? != 0 ]] then
        echo "git commit failed while taking the crontabs backup"| mail -s "[CRONTABS_BACKUP] git commit failed" $EMAIL
        exit 1
    fi

    git push

    if [[ $? != 0 ]] then
         echo "git push failed while taking crontabs backup"|mail -s "[CRONTABS_BACKUP]  git push failed" $EMAIL
        exit 1
    fi
    echo "Cronjobs backup from production is taken" | mail -s "[CRONTABS_BACKUP] Committedi the crons from producion servers" $EMAIL
}


if [ "$VCS" ==  "svn" ]; then
    echo "Using svn as  vcs provider"
    svn_backup
else
    echo "Using git as vcs provider"
    git_backup
fi


