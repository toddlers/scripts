#!/bin/sh
#
# file:         tolower
# purpose:      renames file to lower case
#

if [ $# != 1 ] || [ $1 = "--help" ]
then
    echo "Usage: tolower <filename>"
    exit 0
fi

d=`dirname $1`
f=`basename $1`
lf=`echo $f | tr A-Z a-z`
if [ -r $d/$f ]
then
    if [ $f != $lf ]
    then
        (mv $d/$f $d/$lf) && echo "$f -> $lf"
    else
        echo $d/$f is already lowercase
    fi
else
    echo $d/$f not exists
fi


