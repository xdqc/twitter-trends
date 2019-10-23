import numpy as np
from wordcloud import WordCloud
from random import shuffle 
from PIL import Image, ImageDraw, ImageFont
import math, sys, os

filename = '' if len(sys.argv)<2 else sys.argv[1]

def generate_word_cloud(file):
    words = []
    with open(file, 'r', encoding='utf-8') as r:
        for line in r.readlines():
            if '.' in line.split(',')[1]:
                freq = float(line.split(',')[1])
                timeu = 2 if line.split(',')[2] == '0' else float(line.split(',')[2])
                word, count = line.split(',')[0].capitalize(), int(freq*1e5*timeu)
                if len(word) > 1:
                    words.extend([word]*count)
    shuffle(words)

    # use mask
    mask_img =  np.array(generate_mask_image(file))

    wc = WordCloud(relative_scaling=0.1, background_color="white", mask=mask_img, prefer_horizontal=0.8, max_words=400, max_font_size=500, repeat=True,font_path='./font/AmaticSC-Bold.ttf', margin=0)
    
    # generate word cloud
    wc.generate(' '.join(words))

    # store to file
    wc.to_file('./word-cloud/'+file.split('/')[-1].split('.')[0]+'.png')

def generate_mask_image(text):
    text = text.split('-')[-1].split('.')[0][2:]
    text = text[4:]+text[2:4]+text[:2]
    img = Image.new('RGB', (1280, 360), color = (255, 255, 255))

    fnt = ImageFont.truetype('./font/gabo.otf', 400)
    d = ImageDraw.Draw(img)
    d.text((0,-70), text, font=fnt, fill=(0, 0, 0))

    return img


def main():
    if filename:
        generate_word_cloud(filename)
    else:
        directory = './tweets-trend/'
        for file in os.listdir(directory):
            if os.path.isfile(directory+file):
                generate_word_cloud(directory+file)

main()