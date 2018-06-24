#!/bin/bash

directory="./collected_tweets"
if [ ! -z "$1" ];then
    directory=$1
fi

today=$(date +%Y%m%d)

for file in ${directory}/*; do
    prev=`echo $file | sed 's/.*-\([0-9]\{8\}\)-.*/\1/g'`
    echo $prev
    if [ "$prev" -lt "$today" ];then
        comb=`echo $file | sed 's/\(-[0-9]\{8\}\).+/\1/g'`
        cat $comb* >> $comb.json
        rm -f $comb-*
    fi
done

python tweetText.py
python trendWords.py
