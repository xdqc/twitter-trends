## What’s hot on twitter right now?

![trend-20180927][wordcloud]

[wordcloud]: https://raw.githubusercontent.com/xdqc/tweet-trend-everyday/master/word-cloud/trend-20180927.png "trend-20180927"

## Text mining procedures:

1. Get raw tweets via twitter API (English tweets only)

2. Extract text parts as well as hashtags of tweets for each day.

3. Count unigram tokens (i.e. one English word as a token) of tweet texts, calculate the frequency of each token, order by frequency high to low.

4. Compare today’s model and yesterday’s model. If today’s frequency of one particular token is higher than 2 times of yesterday’s frequency, which mean the hottiness of that word increased by more than 100%, we record that token. (Only care about tokens with frequency >10^-5 )

5. Generate word cloud images.
