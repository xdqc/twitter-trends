"""Get random tweets from Twitter

"""
import re
import json
import yaml
import datetime
from twitter import OAuth, TwitterStream


with open("config.yml", 'r') as stream:
    try:
        c = yaml.load(stream)
        oauth = OAuth(c['Access_Key'], c['Access_Secret'],
                      c['API_Key'], c['API_Secret'])
    except yaml.YAMLError as exc:
        print(exc)

tweet_count = 100000

twitter_stream = TwitterStream(auth=oauth)

iterator = twitter_stream.statuses.sample()

RE_EMOJI = re.compile(
    '[\u0020\u0022\u005C\u0080-\U0010FFFF]', flags=re.UNICODE)

# \uFFFD is the replace char, need to be substituted back to "'"
TR_QUOTE = str.maketrans("’�", "''")

for tweet in iterator:
    if 'lang' in tweet and 'en' in tweet['lang']:
        if 'text' in tweet:
            # Extended Tweets was introduced when 280-character Tweets were launched in November 2017.
            text = tweet['extended_tweet']['full_text'] if 'extended_tweet' in tweet else tweet['text']
            text = re.sub(r'( \w){4,}[^\w]', ' ', text)  # Compress 'Y E L L '
            words = [RE_EMOJI.sub(r'', w.translate(TR_QUOTE)) for w in text.split() if re.match(r'\w', w) and not any(
                x in w for x in ['@', '#', '&', '\\', '://', 'https:', 'http:']) and not w.upper() == 'RT']

            time = int(datetime.datetime.strptime(
                tweet['created_at'], '%a %b %d %H:%M:%S %z %Y').timestamp())
            hashtags = json.dumps(tweet['extended_tweet']['entities']['hashtags']
                                  ) if 'extended_tweet' in tweet else json.dumps(tweet['entities']['hashtags'])

            coordinates = json.dumps(tweet['coordinates'])
            place = json.dumps(tweet['place'])
            user = json.dumps(tweet['user'])

            print('{"created_at":' + str(time) + ', "text":"' + ' '.join(words) + '", "hashtags":' +
                  hashtags + ', "coordinates":' + coordinates + ', "place":' + place + '}')

            # print(json.dumps(tweet))
            tweet_count -= 1
            if tweet_count <= 0 or -1 < time % 3600 < 2:
                break

