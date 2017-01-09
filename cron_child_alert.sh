#!/usr/bin/env bash
# this will print the long running child of cron
# usage : cron_child_alert.sh <threshold in sec>

threshold=$1
cronpid="$(pgrep -P 1 cron)"
cronjobs="$(pstree -p $cronpid | sed 's/)/)\n/g' | grep -v cron | grep -wo '[0-9]*')"
for job in ${cronjobs[@]}
do
  ptime="$(ps -p $job -o etimes,cmd --no-headers|grep bash|awk '{print $1}'|tr -d "\n")"
  if [ $ptime > $threshold ]
  then
    pcron="$(ps -p "$job" -o etimes,cmd --no-headers|grep bash)"
    echo "Long Running Cron Child $pcron"
  fi
done
