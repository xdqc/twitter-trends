package main

import (
	"flag"
	"log"
	"time"

	tweet "github.com/xdqc/dsm-assgn1-tweet"
)

var (
	mode        int
	approach    int
	counterSize int
	directory   string
	output      string
	chinese     int
)

func init() {
	flag.IntVar(&mode, "m", 0, "Two modes: -m 0 static counter; -m <run time in minutes> stream counter")
	flag.IntVar(&approach, "a", 0, "Two approaches: -a 1 run individual counter parallel (should be the requirment); -A 2 for each #tag do TZ/Word counters (my previous understanding of the assignment)")
	flag.IntVar(&counterSize, "n", 0, "The size of space-saving counter; -n 0 to use memery-hungry counter instead")
	flag.StringVar(&directory, "i", "", "The input *.json files directory")
	flag.StringVar(&output, "o", "result.csv", "The output filename for the result")
	flag.IntVar(&chinese, "c", 0, "-c 1 only count Chinese tweets, use *Jieba to segregate words, rather than space")
	flag.Parse()
}

func main() {
	tweet.GetConfig()

	start := time.Now()
	if !(approach == 1 || approach == 2) {
		log.Fatalln("Wrong #Approach. -h for help")
	}

	onlyCountChinese := false
	if chinese == 1 {
		onlyCountChinese = true
	}

	if mode == 0 {
		tweet.Run(approach, directory, counterSize, output, onlyCountChinese)
	} else if mode > 0 {
		tweet.RunStream(approach, counterSize, mode, onlyCountChinese)
	}

	log.Printf("Job done... Process time: %.2f s\n", time.Now().Sub(start).Seconds())

}
