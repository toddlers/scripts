#!/bin/bash

# This pre-commit hook checks following things:
#   1. Empty Commit messages
#   2. Commit message length
#   3. Check for proper jira/bug commit id

REPOS="$1"
TXN="$2"

# set the svnlook command path

SVNLOOK=/usr/bin/svnlook

# Get the log message

LOGMSG=$($SVNLOOK log -t "$TXN" "$REPOS" | grep [a-zA-Z0-9] | wc -c)

# Make sure that the log message contains some text.

if [ "$LOGMSG" = 0 ]; then
echo "Empty log messages are not allowed. Please provide a proper log message" 1>&2
exit 1

#Check that commit message is more than 50 characters long

elif [ "$LOGMSG" -lt 50 ];then
echo -e "Please provide a meaningful comment when committing changes." 1>&2
exit 1

# Check log message for proper task/bug identification
elif [ -x ${REPOS}/hooks/check_log_message.sh ]; then
${REPOS}/hooks/check_log_message.sh "${REPOS}" "${TXN}" 1>&2 || exit 1
else
exit 0
fi
