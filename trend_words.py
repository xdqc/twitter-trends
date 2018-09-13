import os, sys, operator, math
from collections import Counter

directory = './tweets-model/' if len(sys.argv)<2 else sys.argv[1]

days = []

for file in sorted(os.listdir(directory)):
    if file.endswith('.csv'):
        days.append(file)


for i in range(len(days)-1):
    prevModel = {}
    currModel = {}
    hotWords = {}

    with open(directory+days[i], 'r') as f:
        tokens = [t for t in f.readlines() if t]
        for token in tokens:
            word, prob = token.split(',')[0], float(token.split(',')[1])
            prevModel[word] = prob

    with open(directory+days[i+1], 'r') as f:
        tokens = [t for t in f.readlines() if t]
        for token in tokens:
            word, prob = token.split(',')[0], float(token.split(',')[1])
            if prob > 1e-5:
                currModel[word] = prob

    for token in currModel.keys():
        if token in prevModel:
            rate = currModel[token] / prevModel[token]
            if rate > 2.0:
                hotWords[token] = (currModel[token], math.log2(rate))
        else:
            hotWords[token] = (currModel[token], 0)

    with open('tweets-trend/trend-'+days[i+1].split('-')[1]+'.csv', 'w') as f:
        f.write('word, today\'s frequency, log2(freq.today/freq.ystdy)\n')
        [f.write('{0},{1},{2}\n'.format(key, value[0], value[1])) for key, value in hotWords.items()]

    print(days[i+1], len(hotWords), days)
