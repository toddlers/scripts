#!/bin/bash
#Setting the schedular 
find_linux_root_device() {
RDEV=$(mountpoint -d /)

for file in $(find /dev);
do
if [ $(stat --printf="%t:%T" "$file") = $RDEV ];
then
ROOTDEVICE="$file"
break;
fi
done
echo $ROOTDEVICE|cut -d "/" -f3|sed 's/\(.*\)./\1/'
}
while true;
	do
	  result=$(find_linux_root_device)
          echo cfq > /sys/block/$result/queue/scheduler
          echo 256 > /sys/block/$result/queue/nr_requests
	  sleep 3600
done

