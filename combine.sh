#!/bin/bash

directory="./collected_tweets"
if [ ! -z "$1" ];then
    directory=$1
fi

combine() {
    today=$(date +%Y%m%d)
    yesterday=$(date -d 'yesterday' '+%Y%m%d')

    for file in ${directory}/*; do
        prev=`echo $file | sed 's/.*-\([0-9]\{8\}\)-.*/\1/g'`
        echo $prev
        if [ "$prev" -lt "$today" ];then
            comb=`echo $file | sed 's/\(-[0-9]\{8\}\).*/\1/g'`
            cat $comb* >> $comb.json
            rm -f $comb-*
        fi
    done

    python tweet_text.py
    python trend_words.py
    python word_cloud.py ./tweets-trend/trend-${yesterday}.csv

    git add ./tweets-trend/trend-${yesterday}.csv
    git add ./word-cloud/trend-${yesterday}.png
    git commit -m trend-${yesterday}
    git push
}

# do combo task once on 1:00~1:59am each day
while true; do
    if [ $(date +%H) -eq 1 ]; then
        combine
    fi
    sleep 3598
done