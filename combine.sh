#!/bin/bash

directory="./collected_tweets"
if [ ! -z "$1" ];then
    directory=$1
fi

combine() {
    today=$(date +%Y%m%d)
    yesterday=`ls ./collected_tweets | sed -n 2p | cut -c 8-15` #$(date -v -1d '+%Y%m%d')

    count=`ls -1 ${directory}/tweets-${yesterday}-* 2>/dev/null | wc -l`
    if [ $count != 0 ];then 
        cat ${directory}/tweets-${yesterday}-* > ${directory}/tweets-${yesterday}.json
        rm -f ${directory}/tweets-${yesterday}-*
    fi 

    python3 tweet_model.py
    python3 trend_words.py
    python3 word_cloud.py ./tweets-trend-bigram/trend-${yesterday}.csv

    sed -i -e "s/\([0-9]\{8\}\)/${yesterday}/g" readme.md

    # git add ./tweets-trend/trend-${yesterday}.csv
    git add ./word-cloud/trend-${yesterday}.png
    git add ./readme.md
    git commit -m trend-${yesterday}
    #git push

    mv -f ${directory}/tweets-${yesterday}.json ${directory}/archive/tweets-${yesterday}.json
}

# do combo task once on 1:00~1:59am each day
while true; do
    if [ $(date +%H) -eq 1 ]; then
        combine
    fi
    sleep 3598
done
