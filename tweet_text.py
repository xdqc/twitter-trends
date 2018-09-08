import re, os, sys, operator, datetime, difflib
from collections import Counter

directory = './collected_tweets/' if len(sys.argv)<2 else sys.argv[1]


for file in os.listdir(directory):
    sentences = []
    numTweet = 0
    numTweetHashtag = 0

    if len(file) - len(file.replace('-','')) == 1:
        with open(directory+file, 'r') as f:
            tweets = [t for t in f.readlines() if t]
            for tweet in tweets:
                numTweet += 1
                hashtag = re.findall(r'"text": "([^"]+)", "indices"', tweet)
                text = re.findall(r'"text":"([^"]*)", "hashtags"', tweet)[0]
                hashtags = []
                if hashtag:
                    numTweetHashtag += 1
                    hashtags.extend(h.upper() for h in hashtag if h.find('\\')<0)   # ignore non-alphabetic hashtag

                text = text.strip()
                if len(text)>10:
                    text = ' '.join(hashtags) + ' ' + text[0].lower() + text[1:]
                    sentences.append(text)

        print(file,'total tweets:',numTweet, '; tweets with hashtag:', numTweetHashtag, numTweetHashtag/numTweet)
        print(file,'total sentences:', len(sentences), '; unique sentences:', len(set(sentences)))

        freq = Counter(sentences)
        freq = sorted(freq.items(), key=lambda x: (x[0],len(x[0])), reverse=False)
        print('sorted')

        sentences = [x[0] for x in freq]

        # remove similar text
        for i in range(len(sentences)-1):
            # diffstr = ''.join(difflib.ndiff(sentences[i], sentences[i+1]))
            # pluses = len(re.findall(r'[+]', diffstr))
            # minues = len(re.findall(r'[-]', diffstr))
            # if minues/len(sentences[i]) < 0.15:
            #     # if  minues/len(sentences[i]) > 0.14:
            #     #     print('+',pluses,'\t-', minues, '\tl', len(sentences[i]))
            #     #     print(sentences[i])
            #     #     print(sentences[i+1])
            if sentences[i] in sentences[i+1] or sentences[i+1] in sentences[i]:
                sentences[i] = ''

        sentences = [x for x in sentences if x]

        print(file, 'distinct unique sentences:',len(sentences))
        words = []

        # count all words in tweet text
        for text in sentences:
            words.extend(w.strip('\'').lower() for w in re.findall(r'[\da-zA-Z\']+', text) if w.replace('\'',''))

        totalnumberofwords = len(words)
        print(file, 'total number of words:', totalnumberofwords)
        freq = Counter(words)
        print(file, 'unique words', len(freq))
        freq = sorted(freq.items(), key=operator.itemgetter(1), reverse=True)

        outfile = 'tweets-model/' + file.split('.')[0] + '-model.csv'
        with open(outfile, 'w') as f:
            [f.write('{0},{1}\n'.format(item[0], item[1]/totalnumberofwords)) for item in freq]


