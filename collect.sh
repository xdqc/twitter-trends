#!/bin/bash

directory="."
if [ ! -z "$1" ]
then
    directory=$1
fi

for (( i = 1 ; ; i++)) ; do 
    fileName=${directory}/tweets-$(date +%Y%m%d-%H%M).json
    python just_tweets.py > $fileName
    # gzip -9 $fileName &
done
