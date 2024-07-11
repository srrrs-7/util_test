#!/bin/sh

# set args
arg1=$1;
arg2=$2;
arg3=$3;
length=$4;

rm -f ./data.txt;

for ((i = 1; i <= $length; i++));
do
    if [ $i == $length ]
    then
        echo "('$arg1$i', $i, '$arg3')" >> ./data.txt;
        echo ";" >> ./data.txt;
    else
        echo "('$arg1$i', $i, '$arg3')," >> ./data.txt;
    fi
done;

echo "create database data";
