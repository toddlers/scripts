#!/usr/bin/env bash


path=$1

if [[ -z $path ]];then
  echo "please provide the path"
  exit 1
fi

if [[ ! -d $path ]];then
  echo "dir doesnt exists"
  exit 1
fi

index=0

for f in $(find $path -maxdepth 1);do
  if [[ ! -d $f ]];then
    if [[ -L $f ]];then
      fname=$(readlink -f $f)
      files[$index]=$fname
      (( index++ ))
    else
    files[$index]=$f
    (( index++ ))
  fi
fi
done


echo ${files[@]} |awk 'BEGIN{RS=" ";} {print $1}' | awk -F "/" '{print $NF}'|awk -F "." '{print $1}'|sort|uniq
