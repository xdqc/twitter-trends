package tweet

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/yanyiwu/gojieba"

	ss "github.com/xdqc/dsm-assgn1-tweet/spacesaving"
)

//Use mutex to solve concurrent map read and map write problem
var (
	mutex                              sync.Mutex
	mutex2                             sync.Mutex
	mutex3                             sync.Mutex
	JB                                 *gojieba.Jieba
	numTopHeavyHitterThatHasSubcounter = 1000
)

//Run - batch count for saved tweets
func Run(approach int, dir string, counterSize int, outFile string, chinese bool) {
	if chinese {
		JB = gojieba.NewJieba()
	}

	/* Appproach 1 - count hashtag, hashtag&timezone, hashtag&word parallelly */
	hstgCounter := ss.NewCounter(counterSize, false)
	timeZoneHstgCounter := ss.NewCounter(counterSize, false)
	wordHstgCounter := ss.NewCounter(counterSize, false)

	/* Approach 2 - count timezone and word associated with hashtag (under each #tag) */
	hashtagAssociateCounter := ss.NewCounter(counterSize, true) // used for approach2

	var wg sync.WaitGroup
	for _, file := range filesInDirectory(dir) {
		//process tweet files concurrently
		go processTweetFile(approach, file, hstgCounter, timeZoneHstgCounter, wordHstgCounter, hashtagAssociateCounter, chinese, &wg)
	}
	wg.Wait()

	// Output results
	if approach == 1 {
		outputToCSV1(hstgCounter, timeZoneHstgCounter, wordHstgCounter, outFile)
	} else if approach == 2 {
		outputToCSV2(hashtagAssociateCounter, outFile)
	}
}

// Process a tweet file.
// The first three counters args for approach #1, the last counter arg for approach #2
func processTweetFile(approach int, filename string, hstgCtr ss.Counter, tzCtr ss.Counter, wdCtr ss.Counter, hashtagCounter ss.Counter, chinese bool, wg *sync.WaitGroup) {
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
		words := make([]string, 0)

		if chinese {
			if strings.Index(t.Lang, "zh") < 0 {
				continue
			} else {
				useHMM := true
				words = JB.Cut(t.Text, useHMM)
			}
		} else {
			words = strings.Split(t.Text, " ")
		}

		for _, hashtag := range hashtags {
			if approach == 1 {
				//Approach1: count hashtag, hashtag&timezone, hashtag&word parallelly
				countParallel(hashtag.Text, tz, words, wg, hstgCtr, tzCtr, wdCtr)
			} else if approach == 2 {
				//Approach2: count timezone and word under each hashtag
				countPerHashtagAssociate(hashtag.Text, tz, words, hashtagCounter)
			}

		}
	}
}

// Approach1: count pararrell
func countParallel(hashtag string, tz string, words []string, wg *sync.WaitGroup, counters ...ss.Counter) {
	go func() {
		wg.Add(1)
		defer wg.Done()
		mutex.Lock()
		counters[0].Hit(hashtag)
		mutex.Unlock()
	}()
	go func() {
		wg.Add(1)
		defer wg.Done()
		mutex2.Lock()
		counters[1].Hit(hashtag + " * " + tz)
		mutex2.Unlock()
	}()
	go func() {
		wg.Add(1)
		defer wg.Done()
		for _, word := range words {
			mutex3.Lock()
			counters[2].Hit(hashtag + " * " + word)
			mutex3.Unlock()
		}
	}()
}

// Approach2: count timezone and word under each hashtag
func countPerHashtagAssociate(hashtag string, timezone string, words []string, counter ss.Counter) {
	/* Tried to put smaller mutex locked `sync block` in spacesaving package,
		however, prune to cause DEADLOCK, and very hard to debug.
	Just put a big mutex locked block here, may be detrimental to efficiency though. */
	mutex.Lock()

	//count hashtags
	counter.Hit(hashtag)

	if tzCounter := counter.GetSubCounter(hashtag, 0); tzCounter != nil {
		//count timezone associated with the hashtag, use the 0-th subcounter of buckets of hashtagCouter as Timezon counter
		counter.GetSubCounter(hashtag, 0).Hit(timezone)

		//count word associated with the hashtag, use the 1-th subcounter of buckets of hashtagCouter as Word counter
		for _, word := range words[0:1] {
			counter.GetSubCounter(hashtag, 1).Hit(word)
		}
	}
	mutex.Unlock()
}
