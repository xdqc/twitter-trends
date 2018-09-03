package tweet

import (
	"log"
	"net/url"
	"os"
	"os/signal"
	"regexp"
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
func RunStream(approach int, counterSize int, runTimeMinuts int, language string) {
	JB = gojieba.NewJieba()
	resep = regexp.MustCompile("[\n\\p{Z}]")              //tweet tex separator
	repun = regexp.MustCompile("[\n .⠀“ˆ^`｀:،|!　\\p{P}]") //tweet punctuation
	relang := regexp.MustCompile(`(\S{2})`)               //split language arg, every 2 runes
	languages := make([]string, 0)
	for _, lang := range relang.FindAllStringSubmatch(language, -1) {
		languages = append(languages, lang[0])
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
	signal.Notify(sig, syscall.SIGINT) //CTRL+T

	for {
		select {
		case v := <-stream.C:
			tweet, ok := v.(anaconda.Tweet)
			if !ok {
				// Skip bad data
				continue
			}

			if language != "" {
				hasLang := false
				for _, l := range languages {
					if strings.Contains(tweet.Lang, l) {
						hasLang = true
						break
					}
				}
				if !hasLang {
					continue
				}
			}

			go processTweetStream(tweet, approach, hstgCounter, timezoneHstgCounter, wordHstgCounter, hashtagAssociateCounter)

		case <-stop:
			stream.Stop()
			log.Println("Time up")

			//output results to file
			filename := "stream_result/" + strings.Replace(time.Now().Format(time.RFC3339), ":", "", -1) + "_" + strconv.Itoa(runTimeMinuts) + language + ".csv"
			if approach == 1 {
				outputToCSV1(hstgCounter, timezoneHstgCounter, wordHstgCounter, filename)
			} else if approach == 2 {
				outputToCSV2(hashtagAssociateCounter, filename)
			}
			return

		case <-sig:
			log.Printf("Stream counter has been running for %v\n", time.Since(start))

			//output sketchy results to file when user press ctrl+t
			filename := "stream_result/" + strings.Replace(time.Now().Format(time.Stamp), " ", "-", -1) + "_" + strconv.Itoa(int(time.Since(start).Minutes())) + language + "_T.csv"
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
func processTweetStream(t anaconda.Tweet, approach int, counters ...ss.Counter) {

	// // Only count tweet with content
	// if len(t.Text) <= 2 {
	// 	return
	// }

	hashtags := t.Entities.Hashtags
	tz := t.User.TimeZone
	words := make([]string, 0)

	if strings.Index(t.Lang, "zh") >= 0 {
		tokens := JB.Cut(t.Text, true)
		for _, word := range tokens {
			words = append(words, repun.ReplaceAllString(word, ""))
		}
	} else {
		tokens := resep.Split(t.Text, -1)
		for _, word := range tokens {
			words = append(words, repun.ReplaceAllString(word, ""))
		}
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
