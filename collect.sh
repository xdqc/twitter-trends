#!/bin/bash

directory="./collected_tweets"
if [ ! -z "$1" ]
then
    directory=$1
fi

for (( i = 1 ; ; i++)) ; do 
    fileName=${directory}/tweets-$(date +%Y%m%d-%H%M).json
    python just_full_tweets.py > $fileName
    # gzip -9 $fileName &
done
