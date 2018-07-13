import matplotlib.pyplot as plt
from wordcloud import WordCloud
from random import shuffle 
import math, sys, os

directory = './tweets-trend/' if len(sys.argv)<2 else sys.argv[1]

for file in os.listdir(directory):
    words = []
    corpus = {}
    if os.path.isfile(directory+file):
        with open(directory+file, 'r', encoding='utf-8') as r:
            for line in r.readlines():
                if '.' in line.split(',')[1]:
                    word, count = line.split(',')[0].capitalize(), int(math.sqrt(float(line.split(',')[1])*1e5*float(line.split(',')[2])))
                    if len(word) > 1:
                        words.extend([word]*count)
        shuffle(words)
        my_wordcloud = WordCloud(relative_scaling=0.5, width=1080, height=720, max_words=100, min_font_size=20, font_path='./font/AmaticSC-Bold.ttf', margin=0, background_color='white').generate(' '.join(words))
        plt.imshow(my_wordcloud)
        plt.axis("off")
        plt.savefig('./word-cloud/'+file.split('.')[0]+'.png')
