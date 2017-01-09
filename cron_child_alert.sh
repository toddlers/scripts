#!/usr/bin/env bash
# this will print the long running child of cron
# usage : cron_child_alert.sh <threshold in sec>


threshold=$1
status=0
cronpid="$(pgrep -P 1 cron)"
cronjobs="$(pstree -p $cronpid | sed 's/)/)\n/g' | grep -v cron | grep -wo '[0-9]*')"

for job in ${cronjobs[@]}
do
	 #TODO : read the job starttime from /proc
	 ptime="$(ps -p $job -o etimes,cmd --no-headers|grep '/bin/sh'|awk '{print $1}'|tr -d "\n")"
	if [ $ptime > $threshold ]
	then
		#pcron="$(ps -p "$job" -o cmd --no-headers|grep '/bin/sh')"
		pcommand="$(cat /proc/$job/cmdline)"
        status=1
		echo "Long Running Cron Child $job : $pcommand"
	fi
done
exit $status
