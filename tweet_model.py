import re, os, sys, operator, datetime, difflib
from collections import Counter

directory = './collected_tweets/' if len(sys.argv)<2 else sys.argv[1]

unigram_model_dir = './tweets-model/'
bigram_model_dir = './tweets-model-bigram/'

def get_sentences(file):
    numTweet = 0
    numTweetHashtag = 0
    sentences = []
    with open(directory+file, 'r') as f:
        tweets = [t for t in f.readlines() if t]
        for tweet in tweets:
            numTweet += 1
            hashtag = re.findall(r'"text": "([^"]+)", "indices"', tweet)
            hashtags = []
            if hashtag:
                numTweetHashtag += 1
                hashtags.extend(h.upper() for h in hashtag if h.find('\\')<0)   # ignore non-alphabetic hashtag

            textPossible = re.findall(r'"text":"([^"]*)", "hashtags"', tweet)
            if textPossible:
                text = textPossible[0].strip()
                if len(text)>25 or len(text.split())>5:
                    text = ' '.join(hashtags) + ' ' + text[0].lower() + text[1:]
                    sentences.append(text)

    print(file,'total tweets:',numTweet, '; tweets with hashtag:', numTweetHashtag, numTweetHashtag/numTweet)
    print(file,'total sentences:', len(sentences), '; unique sentences:', len(set(sentences)))

    freq = Counter(sentences)
    freq = sorted(freq.items(), key=lambda x: (x[0],len(x[0])), reverse=False) # sort alphabetic then length

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
    return sentences


def build_model(file, tokens, gramliteral):
    numtokens = len(tokens)
    print(file, 'number', gramliteral, numtokens)
    freq = Counter(tokens)
    print(file, 'unique', gramliteral, len(freq))
    freq = sorted(freq.items(), key=operator.itemgetter(1), reverse=True)
    if gramliteral == 'unigrams':
        outfile = 'tweets-model/' + file.split('.')[0] + '-model-'+ str(numtokens) + '.csv'
        with open(outfile, 'w') as f:
            [f.write('{0},{1}\n'.format(item[0], item[1]/numtokens)) for item in freq]
    if gramliteral == 'bigrams':
        outfile = 'tweets-model-bigram/' + file.split('.')[0] + '-bigramodel-'+ str(numtokens) + '.csv'
        with open(outfile, 'w') as f:
            [f.write('{0} {1},{2}\n'.format(item[0][0], item[0][1], item[1]/numtokens)) for item in freq if item[1] > 1]

def build_models():
    for file in os.listdir(directory):
        # only process on concatenated json
        if len(file) - len(file.replace('-','')) == 1:
            sentences = get_sentences(file)
            unigrams, bigrams = [], []
            # count all words in tweet text
            for text in sentences:
                words_in_a_tweet = [w.strip('\'').lower() for w in re.findall(r'[\da-zA-Z\']+', text) if w.replace('\'','')]
                unigrams.extend(words_in_a_tweet)
                bigrams_in_a_tweet = set([b for b in zip(words_in_a_tweet[:-1], words_in_a_tweet[1:])])
                bigrams.extend(bigrams_in_a_tweet)
            # build models 
            build_model(file, unigrams, 'unigrams')
            build_model(file, bigrams, 'bigrams')


build_models()
