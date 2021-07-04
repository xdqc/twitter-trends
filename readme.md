## What’s hot on twitter?

![trend-20210430][wordcloud]

[wordcloud]: https://raw.githubusercontent.com/xdqc/twitter-trend-daily/master/word-cloud/trend-2021/trend-202104/trend-20210430.png "trend-20210430"

### Check trends in history
[2018](https://github.com/xdqc/twitter-trend-daily/tree/master/word-cloud/trend-2018)
[2019](https://github.com/xdqc/twitter-trend-daily/tree/master/word-cloud/trend-2019)
[2020](https://github.com/xdqc/twitter-trend-daily/tree/master/word-cloud/trend-2020)
[2021](https://github.com/xdqc/twitter-trend-daily/tree/master/word-cloud/trend-2021)

### Data stream text mining procedure:

1. **Fetch** tweet stream via twitter API

2. **Count** tweets and hashtags with space saving counters

3. **Build Model** using unigram and bigram tokens

4. **Find Trend** comparing current model with previous models

5. Generate **Word Cloud** images.

