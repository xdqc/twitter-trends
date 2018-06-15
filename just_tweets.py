"""Get random tweets from Twitter

"""
import json, re, datetime
from twitter import OAuth, TwitterStream
import yaml

with open("config.yml", 'r') as stream:
    try:
        c = yaml.load(stream)
        oauth = OAuth(c['Access_Key'], c['Access_Secret'], c['API_Key'], c['API_Secret'])
    except yaml.YAMLError as exc:
        print(exc)

tweet_count = 100000

twitter_stream = TwitterStream(auth=oauth)

iterator = twitter_stream.statuses.sample()

RE_EMOJI = re.compile('[\u0020\u0022\u005C\u0080-\U0010FFFF]', flags=re.UNICODE)

for tweet in iterator:
    if 'lang' in tweet and 'en' in tweet['lang']:
        if 'text' in tweet: 
            text = tweet['extended_tweet']['full_text'] if 'extended_tweet' in tweet else tweet['text']     # Extended Tweets was introduced when 280-character Tweets were launched in November 2017.
            text = re.sub(r'( \w){4,}[^\w]', ' ', text)
            words = [RE_EMOJI.sub(r'', w.replace('�','\'')) for w in text.split() if re.match(r'\w', w) and not any (x in w for x in ['@','#','&','\\','://','https:','http:']) and not w.upper() =='RT']     #\uFFFD is the replace char, need to be substituted back to \'
            if len(' '.join(words).split()) >= 10:
                time = int(datetime.datetime.strptime(tweet['created_at'],'%a %b %d %H:%M:%S %z %Y').timestamp())
                hashtags = json.dumps(tweet['extended_tweet']['entities']['hashtags']) if 'extended_tweet' in tweet else json.dumps(tweet['entities']['hashtags'])
                print('{"created_at":' + str(time) + ', "text":"'+ ' '.join(words) +'", "hashtags":'+ hashtags + '}')
                # print(json.dumps(tweet))
                tweet_count -= 1
                if tweet_count <= 0:
                    break
