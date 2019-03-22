## What’s hot on twitter right now?

![trend-20190320][wordcloud]

[wordcloud]: https://raw.githubusercontent.com/xdqc/tweet-trend-everyday/master/word-cloud/trend-20190320.png "trend-20190320"

## Text mining procedures:

1. **Fetch** raw tweet objects via twitter API (English tweets only).

2. **Extract** texts + hashtags parts of tweets for each day.

3. **Build Model** Count unigram tokens (i.e. one English word as a token) of tweet texts, calculate the occurrence frequency of each token, order by frequency high to low.

4. **Compare** today’s model and yesterday’s model. If today’s frequency of one token is higher than 2 times of yesterday’s frequency, which mean the hottiness of that word increased by more than 100%, record that token. (Only care about tokens with frequency >10^-5 )

5. Generate **Word Cloud** images.
