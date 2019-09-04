#!/bin/bash

directory="./collected_tweets"
if [ ! -z "$1" ];then
    directory=$1
fi

combine() {
    today=$(date +%Y%m%d)
    yesterday=20190826 #$(date -v -1d '+%Y%m%d')

    count=`ls -1 ${directory}/tweets-${yesterday}-* 2>/dev/null | wc -l`
    if [ $count != 0 ];then 
        cat ${directory}/tweets-${yesterday}-* > ${directory}/tweets-${yesterday}.json
        rm -f ${directory}/tweets-${yesterday}-*
    fi 

    python tweet_text.py
    python trend_words.py
    python word_cloud.py ./tweets-trend-bigram/trend-${yesterday}.csv

    sed -i -e "s/\([0-9]\{8\}\)/${yesterday}/g" readme.md

    git add ./tweets-trend/trend-${yesterday}.csv
    git add ./word-cloud/trend-${yesterday}.png
    git add ./readme.md
    git commit -m trend-${yesterday}
    git push

    mv -f ${directory}/tweets-${yesterday}.json ${directory}/archieve/tweets-${yesterday}.json
}

# do combo task once on 1:00~1:59am each day
while true; do
    if [ $(date +%H) -eq 17 ]; then
        combine
    fi
    sleep 3598
done