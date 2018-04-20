package main

import (
	"flag"
	"log"
	"time"

	tweet "github.com/xdqc/dsm-assgn1-tweet"
)

var (
	approach    int
	counterSize int
	directory   string
	output      string
)

func init() {
	flag.IntVar(&approach, "a", 0, "Two approaches: -a 1 run individual counter parallel (should be the requirment); -A 2 associate TZ/Word with each #tag (my previous understanding of the assignment)")
	flag.IntVar(&counterSize, "n", 0, "The size of space-saving counter; -n 0 to use memery-hungry counter instead")
	flag.StringVar(&directory, "i", "", "The directory of input tweet json files")
	flag.StringVar(&output, "o", "result.csv", "The output filename for the result")
	flag.Parse()
}

func main() {
	start := time.Now()
	if !(approach == 1 || approach == 2) {
		log.Fatalln("Wrong #Approach. -h for help")
	}

	tweet.Run(approach, directory, counterSize, output)

	log.Printf("Process time: %.2f s\n", time.Now().Sub(start).Seconds())
}
