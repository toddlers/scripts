#!/bin/bash
#check the md5sum from both, s3 and my local files and remove the local files that are already in amazon s3
datacenter="amazon"
hostname=`hostname`;
path="backup/server245"

s3=`s3cmd ls --list-md5 -H s3://company-backup/company/"$datacenter"/"$hostname"/"$path"/`

s3_list=`echo "$s3"|awk {'print $4" "$5'} | sed 's= .*/= ='`

locally=`md5sum /"$path"/*.gz`;
locally_list=$(echo "$locally" | sed 's= .*/= =');
#echo "$locally_list";

IFS=$'\n'
for i in $locally_list
do
  #echo $i
  locally_hash=`echo $i|awk {'print $1'}`
  locally_file=`echo $i|awk {'print $2'}`

  for j in $s3_list
  do
    s3_hash=$(echo $j|awk {'print $1'}); 
    s3_file=$(echo $j|awk {'print $2'});

    #to avoid empty file when have only hash from folder
    if [[ $s3_hash != "" ]] && [[ $s3_file != "" ]]; then 
      if [[ $s3_hash == $locally_hash ]] && [[ $s3_file == $locally_file ]]; then
        echo "### REMOVING ###";
        echo "$locally_file";
        #rm /"$path"/"$locally_file";
      fi
    fi
  done
done
unset IFS
