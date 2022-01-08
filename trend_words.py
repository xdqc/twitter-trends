import os, sys, operator, math
from collections import Counter

directory_unigram = './tweets-model/'
directory_bigram = './tweets-model-bigram/'


def make_prev_model(directory, days):
    # Use 14 day's summary model as prev
    prevIndex = len(days)-1 if len(days) <= 15 else 15
    prev_model = {}
    prev_words = {}
    total_words_len = 0
    for i, day in enumerate(days[1:prevIndex], start=1):
        with open(directory+day, 'r') as f:
            # copensate for short days
            daily_words_len = int(day.split('-')[-1].split('.')[0]) + 10000000
            # weight damping for older days
            daily_words_len /= (math.log1p(i) + 1)
            total_words_len += daily_words_len
            for token in f.readlines():
                if token:
                    sptoken = token.split(',')
                    word, prob = sptoken[0], float(sptoken[1])
                    if word in prev_words:
                        prev_words[word] += prob * daily_words_len
                    else:
                        prev_words[word] = prob * daily_words_len
    for word in prev_words:
        prev_model[word] = prev_words[word] / total_words_len
    return prev_model


def make_curr_model(directory, days):
    curr_model = {}
    with open(directory+days[0], 'r') as f:
        for token in [t for t in f.readlines() if t]:
            if token:
                sptoken = token.split(',')
                word, prob = sptoken[0], float(sptoken[1])
                if prob > 1e-5:
                    curr_model[word] = prob
    return curr_model


def make_hot_tokens(prev_model, curr_model):
    hot_tokens = {}
    for word in curr_model.keys():
        if word in prev_model:
            rate = curr_model[word] / prev_model[word]
            if rate > 2.0:
                hot_tokens[word] = (curr_model[word], math.log2(rate))
        else:
            hot_tokens[word] = (curr_model[word], 0)
    return hot_tokens


def output_trend(directory, days, hot_tokens):
    outdir = './tweets-trend/' if directory == './tweets-model/' else './tweets-trend-bigram/'
    with open(outdir+'trend-'+days[0].split('-')[1]+'.csv', 'w') as f:
        f.write('word, today\'s frequency, log2(freq.today/freq.ystdy)\n')
        [f.write('{0},{1},{2}\n'.format(key, value[0], value[1])) for key, value in hot_tokens.items()]
    print(days[0], len(hot_tokens))


"""
Create trend of the day by comparing the model of the day with the previous days models:
1) The frequency should > 1e-5
2) The frequency should at least doubled than previous models
"""
def process_trend(directory):
    days = []
    for file in sorted(os.listdir(directory)):
        if file.endswith('.csv'):
            days.insert(0, file)
    prev_model = make_prev_model(directory, days)
    curr_model = make_curr_model(directory, days)
    hot_tokens = make_hot_tokens(prev_model, curr_model)
    output_trend(directory, days, hot_tokens)


process_trend(directory_unigram)
process_trend(directory_bigram)
