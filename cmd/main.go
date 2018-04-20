package main

import (
	"flag"

	tweet "github.com/xdqc/dsm-assgn1-tweet"
)

var (
	directory   string
	counterSize int
)

func init() {
	flag.StringVar(&directory, "d", "", "The directory of tweet json files")
	flag.IntVar(&counterSize, "n", 0, "The size of spacesaving counter")
	flag.Parse()
}

func main() {
	tweet.Run(directory, counterSize)
}
