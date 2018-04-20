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
var (
	mutex  sync.Mutex
	mutex2 sync.Mutex
	mutex3 sync.Mutex
)

func Run(dir string, counterSize int, outFile string) {
	files := filesInDirectory(dir)
	/* Appproach 1 - count hashtag, hashtag&timezone, hashtag&word parallelly */
	// hashtagCounter := ss.NewCounter(counterSize, false)
	// timeZoneHstgCounter := ss.NewCounter(counterSize, false)
	// wordHstgCounter := ss.NewCounter(counterSize, false)

	/* Approach 2 count timezone and word associated with hashtag */
	hashtagCounter := ss.NewCounter(counterSize, true) // used for approach2

	var wg sync.WaitGroup
	for _, file := range files {
		//process tweet files concurrently
		// go processTweetFile(file, hashtagCounter, timeZoneHstgCounter, wordHstgCounter, &wg)

		go processTweetFile(file, hashtagCounter, &wg) // used for approach2
	}
	wg.Wait()

	for i, elem := range hashtagCounter.GetAll()[:10] {
		fmt.Printf("%6d%20s%10d\n", i, elem.Key, elem.Count)
	}
	outputToCSV(hashtagCounter, outFile)
}

//process a tweet file
func processTweetFile(filename string, hstgCounter *ss.Counter, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

	tweetFile, _ := os.Open(filename)
	defer tweetFile.Close()
	scanner := bufio.NewScanner(tweetFile)

	// for each line in the file, process a tweet
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
			//Approach2: count timezone and word per hashtag
			countPerHashtag(hashtag.Text, tz, words, hstgCounter)

		}
	}
}

//Approach2: count timezone and word per hashtag
func countPerHashtag(hashtag string, timezone string, words []string, counter *ss.Counter) {
	/* Tried to put smaller mutex locked `sync block` in spacesaving package,
		however, prune to cause DEADLOCK, and very hard to debug.
	Just put a big mutex locked block here, may be detrimental to efficiency though. */
	mutex.Lock()
	//count hashtags
	counter.Hit(hashtag)

	//count timezone associated with the hashtag, use the 0-th subcounter of buckets of hashtagCouter as Timezon counter
	counter.GetSubCounter(hashtag, 0).Hit(timezone)

	//count word associated with the hashtag, use the 1-th subcounter of buckets of hashtagCouter as Word counter
	for _, word := range words[0:1] {
		counter.GetSubCounter(hashtag, 1).Hit(word)
	}
	mutex.Unlock()
}

//Walk through all .json files in the directory, return <path><filename>
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
		writer.Write([]string{})

		// Word counter
		writer.Write([]string{"HT_rank", "Hashtag", "W_Rank", "Word", "Count"})
		for k, wd := range counter.GetSubCounter(hstg.Key, 1).GetAll() {
			if wd.Count > 1 {
				values := []string{strconv.Itoa(i), hstg.Key, strconv.Itoa(k), wd.Key, strconv.FormatUint(wd.Count, 10)}
				err := writer.Write(values)
				if err != nil {
					log.Panicln("Cannot write to file: " + err.Error())
				}
			}
		}
		writer.Write([]string{})
	}
}
