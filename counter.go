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

//Use mutex to solve concurrent map read and map write problem
var mutex sync.Mutex

func Run(dir string, counterSize int, outFile string) {
	files := filesInDirectory(dir)
	hashtagCounter := ss.NewCounter(counterSize, true)

	var wg sync.WaitGroup
	wg.Add(len(files))
	for _, file := range files {
		//process tweet files concurrently
		go processTweetFile(file, hashtagCounter, &wg)
	}
	wg.Wait()

	for i, elem := range hashtagCounter.GetAll()[:10] {
		fmt.Printf("%6d%20s%10d\n", i, elem.Key, elem.Count)
	}
	outputToCSV(hashtagCounter, outFile)
}

func processTweetFile(filename string, counter *ss.Counter, wg *sync.WaitGroup) {
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
			/* Tried to put smaller mutex surrounded `sync block` in spacesaving package,
				however, prune to cause DEADLOCK, and very hard to debug.
			Just put a big mutexed block here, may be detrimental to efficiency though. */
			mutex.Lock()

			//count hashtags
			counter.Hit(hashtag.Text)

			//count timezone associated with the hashtag, use the 0-th subcounter of buckets of hashtagCouter as Timezon counter
			counter.GetSubCounter(hashtag.Text, 0).Hit(tz)

			//count word associated with the hashtag, use the 1-th subcounter of buckets of hashtagCouter as Word counter
			for _, word := range words[0:1] {
				counter.GetSubCounter(hashtag.Text, 1).Hit(word)
			}
			mutex.Unlock()
		}
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

//Output the result to csv file
func outputToCSV(counter *ss.Counter, outFile string) {
	file, err := os.Create(outFile)
	if err != nil {
		log.Panicln("Cannot create file: " + err.Error())
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Hashtag counter
	writer.Write([]string{"Rank", "Hashtag", "Count", "Uniq_Timezones", "Unique_Words"})
	for i, elem := range counter.GetAll() {
		numUniqTZ := len(counter.GetSubCounter(elem.Key, 0).GetAll())
		numUniqWord := len(counter.GetSubCounter(elem.Key, 1).GetAll())
		values := []string{strconv.Itoa(i), elem.Key, strconv.FormatUint(elem.Count, 10), strconv.Itoa(numUniqTZ), strconv.Itoa(numUniqWord)}
		err := writer.Write(values)
		if err != nil {
			log.Panicln("Cannot write to file: " + err.Error())
		}
	}
	writer.Write([]string{})

	// output timezone and words result for first 10 hashtags
	for i, hstg := range counter.GetAll()[:10] {
		// TimeZone counter
		writer.Write([]string{"HT_rank", "Hashtag", "TZ_Rank", "TimeZone", "Count"})
		for j, tz := range counter.GetSubCounter(hstg.Key, 0).GetAll() {
			values := []string{strconv.Itoa(i), hstg.Key, strconv.Itoa(j), tz.Key, strconv.FormatUint(tz.Count, 10)}
			err := writer.Write(values)
			if err != nil {
				log.Panicln("Cannot write to file: " + err.Error())
			}
		}

		// Word counter
		writer.Write([]string{"HT_rank", "Hashtag", "W_Rank", "Word", "Count"})
		for k, wd := range counter.GetSubCounter(hstg.Key, 1).GetAll() {
			values := []string{strconv.Itoa(i), hstg.Key, strconv.Itoa(k), wd.Key, strconv.FormatUint(wd.Count, 10)}
			err := writer.Write(values)
			if err != nil {
				log.Panicln("Cannot write to file: " + err.Error())
			}
		}
	}
}
