## What’s hot on twitter?

![trend-20190315][wordcloud1]

[wordcloud1]: https://raw.githubusercontent.com/xdqc/twitter-trend-daily/master/word-cloud/trend-2019/trend-201903/trend-20190315.png "trend-20190315"

![trend-20200312][wordcloud2]

[wordcloud2]: https://raw.githubusercontent.com/xdqc/twitter-trend-daily/master/word-cloud/trend-2020/trend-202003/trend-20200312.png "trend-20200312"

![trend-20210106][wordcloud3]

[wordcloud3]: https://raw.githubusercontent.com/xdqc/twitter-trend-daily/master/word-cloud/trend-2021/trend-202101/trend-20210106.png "trend-20210106"

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

