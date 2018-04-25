package tweet

import (
	"encoding/csv"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/ChimeraCoder/anaconda"
	ss "github.com/xdqc/dsm-assgn1-tweet/spacesaving"
)

//Walk through all .json files in the directory, return <path><filename>
func filesInDirectory(dir string) (files []string) {
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".json") {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		log.Panicln(err)
	}
	return
}

//Approach1: Output the result to csv file
func outputToCSV1(hstgCtr ss.Counter, tzCtr ss.Counter, wdCtr ss.Counter, outFile string) {
	file, err := os.Create(outFile)
	if err != nil {
		log.Panicln("Cannot create file: " + err.Error())
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Hashtag counter
	mutex.Lock()
	writer.Write([]string{"Rank", "Hashtag", "Count"})
	for i, elem := range hstgCtr.GetAll() {
		values := []string{strconv.Itoa(i), elem.Key, strconv.FormatUint(elem.Count, 10)}
		err := writer.Write(values)
		if err != nil {
			log.Panicln("Cannot write to file: " + err.Error())
		}
	}
	mutex.Unlock()
	writer.Write([]string{})

	// Hashtag&TimeZone counter
	mutex2.Lock()
	writer.Write([]string{"Rank", "Hashtag&TimeZone", "Count"})
	for i, elem := range tzCtr.GetAll() {
		values := []string{strconv.Itoa(i), elem.Key, strconv.FormatUint(elem.Count, 10)}
		err := writer.Write(values)
		if err != nil {
			log.Panicln("Cannot write to file: " + err.Error())
		}
	}
	mutex2.Unlock()
	writer.Write([]string{})

	// Hashtag&Word counter
	mutex3.Lock()
	writer.Write([]string{"Rank", "Hashtag&Word", "Count"})
	for i, elem := range wdCtr.GetAll() {
		values := []string{strconv.Itoa(i), elem.Key, strconv.FormatUint(elem.Count, 10)}
		err := writer.Write(values)
		if err != nil {
			log.Panicln("Cannot write to file: " + err.Error())
		}
	}
	mutex3.Unlock()
	writer.Write([]string{})
}

//Approach2: Output the result to csv file
func outputToCSV2(counter ss.Counter, outFile string) {
	file, err := os.Create(outFile)
	if err != nil {
		log.Panicln("Cannot create file: " + err.Error())
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Hashtag counter
	writer.Write([]string{"Rank", "Hashtag", "Count", "Uniq_Timezones", "Unique_Words"})
	mutex.Lock()
	for i, elem := range counter.GetAll() {
		if counter.GetSubCounter(elem.Key, 0) != nil {
			numUniqTZ := len(counter.GetSubCounter(elem.Key, 0).GetAll())
			numUniqWord := len(counter.GetSubCounter(elem.Key, 1).GetAll())
			values := []string{strconv.Itoa(i), elem.Key, strconv.FormatUint(elem.Count, 10), strconv.Itoa(numUniqTZ), strconv.Itoa(numUniqWord)}
			err := writer.Write(values)
			if err != nil {
				log.Panicln("Cannot write to file: " + err.Error())
			}
		}
	}
	writer.Write([]string{})

	// output timezone and words result for first 1000 hashtags
	for i, hstg := range counter.GetAll() {
		if counter.GetSubCounter(hstg.Key, 0) != nil {

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

			// Word counter, only for count > 1
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
	mutex.Unlock()
}

func writeTweetFile(t anaconda.Tweet, filepath string) {
	bs, err := json.Marshal(t)
	if err != nil {
		log.Panicln("parse tweet error" + err.Error())
	}

	f, err := os.OpenFile(filepath, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	if _, err = f.Write(bs); err != nil {
		panic(err)
	}
}
