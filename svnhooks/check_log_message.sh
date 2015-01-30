#!/bin/bash

REPOS="${1}"
TXN="${2}"

SVNLOOK=/usr/bin/svnlook

LOG_MSG_LINE1=`${SVNLOOK} log -t "${TXN}" "${REPOS}" | head -n1`

if (echo "${LOG_MSG_LINE1}" | egrep '^[a-zA-Z]+[-][1-9][0-9]*[:]?[\s]*.*$' >  /dev/null;) \
|| (echo "${LOG_MSG_LINE1}" | egrep '^[nN][oO][jJ][iI][rR][aA][:]?[\s]*.*$' > /dev/null;) \
|| (echo "${LOG_MSG_LINE1}" | egrep '^\[maven-release-plugin\][\s]*.*$' >  /dev/null;)
then
exit 0
else
echo ""
echo "Your log message does not contain a JIRA Issue identifier (or bad format used)"
echo "The JIRA Issue identifier must be the first item on the first line of the log message."
echo ""
echo "Proper JIRA format:  'AAA-000'"
#echo "JIRA regex: '^[a-zA-Z]+[-][1-9][0-9]*[:]?[\s]*.*$'"
exit 1
fi

