package tweet

import (
	"log"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/ChimeraCoder/anaconda"
	"github.com/yanyiwu/gojieba"

	ss "github.com/xdqc/dsm-assgn1-tweet/spacesaving"
)

//RunStream - Fetch tweets from api and count
func RunStream(approach int, counterSize int, runTimeMinuts int, isChinese bool) {
	if isChinese {
		JB = gojieba.NewJieba()
	}

	// Start Twitter API
	api := anaconda.NewTwitterApiWithCredentials(cfg.AccessKey, cfg.AccessSecret, cfg.APIKey, cfg.APISecret)
	stream := api.PublicStreamSample(url.Values{})
	log.Println("Tweet API working ... will run " + strconv.Itoa(runTimeMinuts) + " minutes.")

	// Create Counters
	hstgCounter := ss.NewCounter(counterSize, false)
	timezoneHstgCounter := ss.NewCounter(counterSize, false)
	wordHstgCounter := ss.NewCounter(counterSize, false)

	hashtagAssociateCounter := ss.NewCounter(counterSize, true) // used for approach2

	// Start timer
	start := time.Now()
	stop := make(chan int)
	go afterTimer(&stop, runTimeMinuts)

	//listen to system signal
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINFO) //CTRL+T

	for {
		select {
		case v := <-stream.C:
			tweet, ok := v.(anaconda.Tweet)
			if !ok {
				// Skip bad data
				continue
			}

			go processTweetStream(tweet, approach, isChinese, hstgCounter, timezoneHstgCounter, wordHstgCounter, hashtagAssociateCounter)

		case <-stop:
			stream.Stop()
			log.Println("Time up")

			//output results to file
			filename := "stream_result/" + strings.Replace(time.Now().Format(time.RFC3339), ":", "", -1) + "_" + strconv.Itoa(runTimeMinuts) + ".csv"
			if approach == 1 {
				outputToCSV1(hstgCounter, timezoneHstgCounter, wordHstgCounter, filename)
			} else if approach == 2 {
				outputToCSV2(hashtagAssociateCounter, filename)
			}
			return

		case <-sig:
			log.Printf("Stream counter has been running for %v\n", time.Since(start))

			//output sketchy results to file when user press ctrl+t
			filename := "stream_result/" + strings.Replace(time.Now().Format(time.RFC3339), ":", "", -1) + "_" + strconv.Itoa(int(time.Since(start).Minutes())) + "_T.csv"
			if approach == 1 {
				go outputToCSV1(hstgCounter, timezoneHstgCounter, wordHstgCounter, filename)
			} else if approach == 2 {
				go outputToCSV2(hashtagAssociateCounter, filename)
			}
		}
	}
}

// Process a tweet in the stream.
// The first three counters for approach #1, the last counter for approach #2
func processTweetStream(t anaconda.Tweet, approach int, chinese bool, counters ...ss.Counter) {

	// // Only count tweet with content
	// if len(t.Text) <= 2 {
	// 	return
	// }

	hashtags := t.Entities.Hashtags
	tz := t.User.TimeZone
	words := make([]string, 0)

	if chinese {
		if strings.Index(t.Lang, "zh") < 0 && strings.Index(t.User.Lang, "zh") < 0 {
			return
		}
		useHMM := true
		words = JB.Cut(t.Text, useHMM)
	} else {
		words = strings.Split(t.Text, " ")
	}

	for _, hashtag := range hashtags {
		if approach == 1 {
			//Approach1: count hashtag, hashtag&timezone, hashtag&word parallelly
			countParallel(hashtag.Text, tz, words, new(sync.WaitGroup), counters[0], counters[1], counters[2])
		} else if approach == 2 {
			//Approach2: count timezone and word under each hashtag
			countPerHashtagAssociate(hashtag.Text, tz, words, counters[3])
		}
	}

	// Also count tweets without hashtag as well
	if len(hashtags) == 0 {
		if approach == 1 {
			//Approach1: count hashtag, hashtag&timezone, hashtag&word parallelly
			countParallel(" ", tz, words, new(sync.WaitGroup), counters[0], counters[1], counters[2])
		} else if approach == 2 {
			//Approach2: count timezone and word under each hashtag
			countPerHashtagAssociate(" ", tz, words, counters[3])
		}
	}
}

func afterTimer(stop *chan int, min int) {
	time.Sleep(time.Duration(min) * time.Minute)
	*stop <- 1
}
