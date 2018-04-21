package tweet

import (
	"net/url"
	"strings"
	"time"

	"github.com/ChimeraCoder/anaconda"
	ss "github.com/xdqc/dsm-assgn1-tweet/spacesaving"
)

//RunStream - Fetch tweets from api and count
func RunStream(approach int, counterSize int, runTimeMinuts int) {
	api := anaconda.NewTwitterApiWithCredentials(cfg.AccessKey, cfg.AccessSecret, cfg.APIKey, cfg.APISecret)
	println(api.Log)
	stream := api.PublicStreamSample(url.Values{})

	hstgCounter := ss.NewCounter(counterSize, false)
	timezoneHstgCounter := ss.NewCounter(counterSize, false)
	wordHstgCounter := ss.NewCounter(counterSize, false)

	hashtagAssociateCounter := ss.NewCounter(counterSize, true) // used for approach2

	// start the timer
	stop := make(chan int)
	go afterTimer(&stop)

	for {
		select {
		case v := <-stream.C:
			tweet, ok := v.(anaconda.Tweet)
			if !ok {
				continue
			}

			go processTweetStream(tweet, approach, hstgCounter, timezoneHstgCounter, wordHstgCounter, hashtagAssociateCounter)

		case <-stop:
			stream.Stop()
			filename := "/stream_result/" + strings.Replace(time.Now().Format(time.RFC3339), ":", "", -1)
			outputToCSV1(hstgCounter, timezoneHstgCounter, wordHstgCounter, filename)
			break
		}
	}
}

func processTweetStream(t anaconda.Tweet, approach int, counters ...*ss.Counter) {

	hashtags := t.Entities.Hashtags
	tz := t.User.TimeZone
	words := strings.Split(t.Text, " ")
	for _, hashtag := range hashtags {
		if approach == 1 {
			//Approach1: count hashtag, hashtag&timezone, hashtag&word parallelly
			go func() {
				mutex.Lock()
				counters[0].Hit(hashtag.Text)
				mutex.Unlock()
			}()
			go func() {
				mutex2.Lock()
				counters[1].Hit(hashtag.Text + " * " + tz)
				mutex2.Unlock()
			}()
			go func() {
				for _, word := range words {
					mutex3.Lock()
					counters[2].Hit(hashtag.Text + " * " + word)
					mutex3.Unlock()
				}
			}()
		}
	}
}

func afterTimer(stop *chan int) {
	time.Sleep(time.Minute * 30)
	*stop <- 1
}
