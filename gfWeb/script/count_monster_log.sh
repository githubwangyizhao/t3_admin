#!/bin/bash

OPTS="cat $1"

if [  "$3" ]; then
    OLD_IFS="$IFS"
    IFS=$2
    arr=($3)
    IFS="$OLD_IFS"

    for s in ${arr[@]}
    do
    OPTS="$OPTS | grep $s"
    done
fi



eval $OPTS |  awk '{print $2}' |awk -F , '{print $4 ,$12}'  | awk -F } '{print $1 }' | sort  | uniq -c | sort -rnk 1
eval $OPTS |  awk '{print $2}' |awk -F , '{print $4 ,$12}'  | awk -F } '{sum += $1};END {print "sum："sum}'
