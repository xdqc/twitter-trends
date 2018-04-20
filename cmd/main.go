package main

import (
	"flag"

	tweet "github.com/xdqc/dsm-assgn1-tweet"
)

var (
	counterSize int
	directory   string
	output      string
)

func init() {
	flag.IntVar(&counterSize, "n", 1000, "The size of spacesaving counter")
	flag.StringVar(&directory, "i", "march18", "The directory of input tweet json files")
	flag.StringVar(&output, "o", "result18.csv", "The output filename for the result")
	flag.Parse()
}

func main() {
	tweet.Run(directory, counterSize, output)
}
