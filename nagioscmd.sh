#!/bin/sh
############################################################################
# Script: cmd.sh
# Description:
#       This script is used to run Nagios commands on the command-line. It passes
#       commands to the Nagios command file by printing the command string and
#       directing it to the command file.
# Usage:
#       ./nagioscmd.sh -h host -s service -c nagios_cmd
#       ./nagioscmd.sh -h db.domain.net -s DB-CHK -c ENABLE_SVC_CHECK
############################################################################
# Define Variables
now=`date +%s`
cmdfile=’/var/log/nagios/rw/nagios.cmd’
hflag=0
sflag=0
cflag=0
# Define usage function
usage(){
echo “Usage: ./cmd.sh -h <host> -s <service> -c <nagios_cmd>”
echo “   -h: FQDN of the host you are running command against”
echo “   -s: Service you are checking (ie: DB-SYNCER)”
echo “   -c: Nagios command options: http://old.nagios.org/developerinfo/externalcommands/commandlist.php&#8221;
echo “”
echo “   Example: ./cmd.sh -h db.domain.net -s DB-CHK -c ENABLE_SVC_CHECK”
echo “”
exit 1
}
if [ $# -ne 6 ]; then
usage
fi
# Extract cmd line parameters
while [ $# -gt 0 ]
do
case “$1″ in
-h) hflag=1;h=$2;;
-s) sflag=1;s=$2;;
-c) cflag=1;c=$2;;
esac
shift
done
if [[ -z $h || -z $c || -z $s ]]
then
echo “ERROR: Missing -h, -c, or -s”
usage
else
printf “[%lu] $c;$h;$s\n” $now > $cmdfile
fi
