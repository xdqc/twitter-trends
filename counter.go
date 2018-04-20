package tweet

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	ss "github.com/xdqc/dsm-assgn1-tweet/spacesaving"
)

func Run(dir string, counterSize int) {
	files := filesInDirectory(dir)
	files = files[0:2]

	//The second param `2` indicates the two sub-counters for timezone and word associated with the hashtag, respectively
	hashtagCounter := ss.NewCounter(counterSize, 2)

	var wg sync.WaitGroup
	wg.Add(len(files))
	for _, file := range files {
		go processTweetFile(file, counterSize, hashtagCounter)
	}
	wg.Wait()

	for i, elem := range hashtagCounter.GetAll() {
		fmt.Printf("%6d%20s%10d\n", i, elem.Key, elem.Count)
	}
}

//Walk through all .json files in the directory, return <path><filenames>
func filesInDirectory(dir string) (files []string) {
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".json") {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	return
}

func processTweetFile(filename string, counterSize int, hashtagCounter *ss.Counter) {
	tweetFile, _ := os.Open(filename)
	defer tweetFile.Close()
	scanner := bufio.NewScanner(tweetFile)

	// for each line of the file, process a tweet
	for scanner.Scan() {
		var t *Tweet
		err := json.Unmarshal(scanner.Bytes(), &t)
		if err != nil {
			log.Println("parse tweet err:" + err.Error())
		}
		hashtags := t.Entities.Hashtags
		tz := t.User.TimeZone
		words := strings.Split(t.Text, " ")

		for _, hashtag := range hashtags {
			//count hashtags
			hashtagCounter.Hit(hashtag.Text)

			//count timezone associated with the hashtag, use the 0-th subcounter of hashtagCouter
			hashtagCounter.GetSubCounter(hashtag.Text, 0).Hit(tz)

			//count word associated with the hashtag, use the 1-th subcounter of hashtagCouter
			for _, word := range words {
				hashtagCounter.GetSubCounter(hashtag.Text, 1).Hit(word)
			}
		}
	}
}
