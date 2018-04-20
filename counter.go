package tweet

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	ss "github.com/xdqc/dsm-assgn1-tweet/spacesaving"
)

func Run(dir string, counterSize int, outFile string) {
	files := filesInDirectory(dir)
	files = files[0:1]

	hashtagCounter := ss.NewCounter(counterSize)

	var wg sync.WaitGroup
	wg.Add(len(files))
	for _, file := range files {
		//process tweet files concurrently
		go func(filename string) {
			defer wg.Done()
			tweetFile, _ := os.Open(filename)
			defer tweetFile.Close()
			scanner := bufio.NewScanner(tweetFile)

			// for each line of the file, process a tweet
			for scanner.Scan() {
				var t *Tweet
				err := json.Unmarshal(scanner.Bytes(), &t)
				if err != nil {
					log.Println("parse tweet err: " + err.Error())
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
		}(file)
	}
	wg.Wait()

	for i, elem := range hashtagCounter.GetAll()[:100] {
		fmt.Printf("%6d%20s%10d\n", i, elem, elem.Count)
	}
	outputToCSV(hashtagCounter, outFile)
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

func outputToCSV(counter *ss.Counter, outFile string) {
	file, err := os.Create(outFile)
	if err != nil {
		log.Panicln("Connet create file: " + err.Error())
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for i, elem := range counter.GetAll() {
		numUniqTZ := len(counter.GetSubCounter(elem.Key, 0).GetAll())
		numUniqWord := len(counter.GetSubCounter(elem.Key, 1).GetAll())
		values := []string{strconv.Itoa(i), elem.Key, strconv.FormatUint(elem.Count, 10), strconv.Itoa(numUniqTZ), strconv.Itoa(numUniqWord)}
		err := writer.Write(values)
		if err != nil {
			log.Panicln("Connet write to file: " + err.Error())
		}
	}
}
