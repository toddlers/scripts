#!/usr/bin/env bash
# this will print the long running child of cron
# usage : cron_child_alert.sh <threshold in sec>

#!/usr/bin/env bash
threshold=$1
status=0
cronpid="$(pgrep -P 1 cron)"
cronjobs="$(pstree -p $cronpid | sed 's/)/)\n/g' | grep -v cron | grep -wo '[0-9]*')"
#echo $cronjobs
for job in ${cronjobs[@]}
do
	#for job in ${cronjobs[@]}; do  ps -p "$job" -o etimes,cmd --no-headers|grep bash|awk '{print $1}'; done
	ptime="$(ps -p $job -o etimes,cmd --no-headers|grep '/bin/sh'|awk '{print $1}'|tr -d "\n")"
	if [ $ptime > $threshold ]
	then
		pcron="$(ps -p "$job" -o cmd --no-headers|grep '/bin/sh')"
        status=1
		echo "Long Running Cron Child $job : $pcron"
	fi
done
exit $statuss
