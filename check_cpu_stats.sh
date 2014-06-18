#!/bin/bash

# CPU usage check

IOSTAT=/usr/bin/iostat

# Nagios return codes
STATE_OK=0
STATE_WARNING=1
STATE_CRITICAL=2
STATE_UNKNOWN=3

# Plugin parameters value if not define
WARNING_THRESHOLD=${WARNING_THRESHOLD:="5"}
CRITICAL_THRESHOLD=${CRITICAL_THRESHOLD:="1"}
INTERVAL_SEC=${INTERVAL_SEC:="1"}
NUM_REPORT=${NUM_REPORT:="3"}

# Plugin variable description
PROGNAME=$(basename $0)

if [ ! -x $IOSTAT ]; then
        echo "UNKNOWN: iostat not found or is not executable by the nagios user."
        exit $STATE_UNKNOWN
fi

print_usage() {
        echo ""
        echo "$PROGNAME $RELEASE - CPU Utilization check script for Nagios"
        echo ""
        echo "Usage: check_cpu_stats.sh -w -c (-i -n)"
        echo ""
        echo "  -w  Warning level in % for cpu iowait"
        echo "  -c  Crical level in % for cpu iowait"
        echo "  -i  Interval in seconds for iostat (default : 1)"
        echo "  -n  Number report for iostat (default : 3)"
        echo "  -h  Show this page"
        echo ""
    echo "Usage: $PROGNAME"
    echo "Usage: $PROGNAME --help"
    echo ""
}

print_help() {
        print_usage
        echo ""
        echo "This plugin will check cpu utilization (user,system,iowait,idle in %)"
        echo ""
        exit 0
}

# Parse parameters
while [ $# -gt 0 ]; do
    case "$1" in
        -h | --help)
            print_help
            exit $STATE_OK
            ;;
        -w | --warning)
                shift
                WARNING_THRESHOLD=$1
                ;;
        -c | --critical)
               shift
                CRITICAL_THRESHOLD=$1
                ;;
        -i | --interval)
               shift
               INTERVAL_SEC=$1
                ;;
        -n | --number)
               shift
               NUM_REPORT=$1
                ;;
        *)  echo "Unknown argument: $1"
            print_usage
            exit $STATE_UNKNOWN
            ;;
        esac
shift
done

CPU_REPORT=`iostat -c $INTERVAL_SEC $NUM_REPORT | sed -e 's/,/./g' | tr -s ' ' ';' | sed '/^$/d' | tail -1`
CPU_REPORT_SECTIONS=`echo ${CPU_REPORT} | grep ';' -o | wc -l`
CPU_USER=`echo $CPU_REPORT | cut -d ";" -f 2`
CPU_NICE=`echo $CPU_REPORT | cut -d ";" -f 3`
CPU_SYSTEM=`echo $CPU_REPORT | cut -d ";" -f 4`
CPU_IOWAIT=`echo $CPU_REPORT | cut -d ";" -f 5`
if [ ${CPU_REPORT_SECTIONS} -ge 6 ]; then
    CPU_STEAL=`echo $CPU_REPORT | cut -d ";" -f 6`
    CPU_IDLE=`echo $CPU_REPORT | cut -d ";" -f 7`
    NAGIOS_DATA="user=${CPU_USER}% system=${CPU_SYSTEM}% iowait=${CPU_IOWAIT}% idle=${CPU_IDLE}% nice=${CPU_NICE}% steal=${CPU_STEAL}%|user=${CPU_USER}%;;;; system=${CPU_SYSTEM}%;;;; iowait=${CPU_IOWAIT}%;;;; idle=${CPU_IDLE}%;;;; nice=${CPU_NICE}%;;;; steal=${CPU_STEAL}%;;;;"
else
    CPU_IDLE=`echo $CPU_REPORT | cut -d ";" -f 6`
    NAGIOS_DATA="user=${CPU_USER}% system=${CPU_SYSTEM}% iowait=${CPU_IOWAIT}% idle=${CPU_IDLE}% nice=${CPU_NICE}%|user=${CPU_USER}%;;;; system=${CPU_SYSTEM}%;;;; iowait=${CPU_IOWAIT}%;;;; idle=${CPU_IDLE}%;;;; nice=${CPU_NICE}%;;;;"
fi
    CPU_ALERTING_METRIC=`echo $CPU_IDLE | cut -d "." -f 1`

if [ ${CPU_ALERTING_METRIC} -le $WARNING_THRESHOLD ] && [ ${CPU_ALERTING_METRIC} -gt $CRITICAL_THRESHOLD ]; then
    echo "CPU STATISTICS WARNING : ${NAGIOS_DATA}"
    exit $STATE_WARNING
elif [ ${CPU_ALERTING_METRIC} -le $CRITICAL_THRESHOLD ]; then
    echo "CPU STATISTICS CRITICAL : ${NAGIOS_DATA}"
    exit $STATE_CRITICAL
else
    echo "CPU STATISTICS OK : ${NAGIOS_DATA}"
    exit $STATE_OK
fi
