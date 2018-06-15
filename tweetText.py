import re, os, sys, operator, datetime, difflib
from collections import Counter

directory = './collected_tweets/' if len(sys.argv)<2 else sys.argv[1]

sentences = []
numTweet = 0
numTweetHashtag = 0

for file in os.listdir(directory):
    if len(file) - len(file.replace('-','')) == 1:
        with open(directory+file, 'r') as f:
            tweets = [t for t in f.readlines() if t]
            for tweet in tweets:
                numTweet += 1
                hashtag = re.findall(r'"text": "([^"]+)", "indices"', tweet)
                text = re.findall(r'"text":"([^"]+)", "hashtags"', tweet)[0]
                hashtags = []
                if hashtag:
                    numTweetHashtag += 1
                    hashtags.extend(h.upper() for h in hashtag if h.find('\\')<0)   # ignore non-alphabetic hashtag
                if text:
                    text = text.strip()
                    text = ' '.join(hashtags) + ' ' + text[0].lower() + text[1:]
                    sentences.append(text)

print('total tweets:',numTweet, '; tweets with hashtag:', numTweetHashtag, numTweetHashtag/numTweet)
print('total sentences:', len(sentences), '; unique sentences:', len(set(sentences)))

freq = Counter(sentences)
freq = sorted(freq.items(), key=lambda x: (x[0],len(x[0])), reverse=False)
print('sorted')

sentences = [x[0] for x in freq]

# remove similar text
for i in range(len(sentences)-1):
    diffstr = ''.join(difflib.ndiff(sentences[i], sentences[i+1]))
    pluses = len(re.findall(r'[+]', diffstr))
    minues = len(re.findall(r'[-]', diffstr))
    if minues/len(sentences[i]) < 0.15:
        if  minues/len(sentences[i]) > 0.1:
            print('+',pluses,'\t-', minues, '\tl', len(sentences[i]))
            print(sentences[i])
            print(sentences[i+1])
        sentences[i] = ''

outfile = 'tweetText.csv'

with open(outfile, 'a') as w:
    for item in sentences:
        if item:
            w.write(item+'\n')