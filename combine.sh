#!/bin/bash

directory="./collected_tweets"
if [ ! -z "$1" ];then
    directory=$1
fi

combine() {
    today=$(date +%Y%m%d)
    yesterday=$(date -v -1d '+%Y%m%d')

    count=`ls -1 ${directory}/tweets-${yesterday}-* 2>/dev/null | wc -l`
    if [ $count != 0 ];then 
        cat ${directory}/tweets-${yesterday}-* > ${directory}/tweets-${yesterday}.json
        rm -f ${directory}/tweets-${yesterday}-*
    fi 

    python tweet_text.py
    python trend_words.py
    python word_cloud.py ./tweets-trend/trend-${yesterday}.csv

    sed -i -E "s/\([0-9]\{8\}\)/${yesterday}/g" readme.md

    git add ./tweets-trend/trend-${yesterday}.csv
    git add ./word-cloud/trend-${yesterday}.png
    git add ./readme.md
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