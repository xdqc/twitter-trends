package spacesaving

import (
	"log"
)

var (
	//use 1000 or 10000 for the number of buckets of the counter
	size int
	//`2` indicates the two sub-counters for timezone and word associated with the hashtag, respectively
	numSubCounters = 2
)

//NewCounter - initialise spacesaving counter, size 0 for memory hungry counter
func NewCounter(s int, isSuperCounter bool) *Counter {
	size = s
	ss := Counter{
		list:    make([]Element, size),
		hash:    make(map[string]uint32, size),
		isSuper: isSuperCounter,
	}
	return &ss
}

//Hit - for every instance in stream
func (ss *Counter) Hit(key string) {
	var (
		idx    uint32
		found  bool
		bucket *Element
	)

	if idx, found = ss.hash[key]; found {
		// exsisting element => increment count
		bucket = &ss.list[idx]
	} else {
		idx = 0
		// Given size 0 for memory-hungry counter
		if size > 0 {
			// Space-Saving
			// new element => replace the first element(lowest count) with new key
			bucket = &ss.list[idx]
			bucket.Key = key
			delete(ss.hash, bucket.Key)
		} else {
			// Memory-Hungry
			// new element => create new Element bucket
			// prepend the new bucket to list (golang doesn't have build-in prepend(). Use append() to append the old list to an one element list)
			bucket = &Element{
				Key:   key,
				Count: 0,
			}
			ss.list = append([]Element{*bucket}, ss.list...)
		}
		ss.hash[key] = idx

		// only create subcounters for supercounter
		if ss.isSuper {
			bucket.subCounters = make([]*Counter, 2)
			for i := 0; i < numSubCounters; i++ {
				bucket.subCounters[i] = NewCounter(size, false)
			}
		}
	}

	// increment count for the bucket
	bucket.Count++

	// sort the list, lower count in front and higher count in end
	for {
		if idx == uint32(len(ss.list))-1 {
			break
		}

		b1 := &ss.list[idx]
		b2 := &ss.list[idx+1]
		//ignore counting ties
		if b1.Count <= b2.Count {
			break
		}

		//switch buckets
		ss.hash[b1.Key] = idx + 1
		ss.hash[b2.Key] = idx
		*b1, *b2 = *b2, *b1
		idx++
	}
}

//GetSubCounter - get subcounter of the bucket with key
func (ss *Counter) GetSubCounter(key string, i int) (subCounter *Counter) {
	if i >= numSubCounters {
		log.Panicln("subcounter index out of bound")
	}
	if elemIdx, found := ss.hash[key]; found {
		elem := ss.list[elemIdx]
		return elem.subCounters[i]
	}
	log.Panicln("call for subcounter under no-exist supercounter's bucket")
	return
}

//GetAll return all elements in the counter
func (ss *Counter) GetAll() (elements []Element) {
	elements = make([]Element, 0, len(ss.hash))
	// output from higt to low count
	for i := len(ss.list) - 1; i >= 0; i-- {
		b := &ss.list[i]
		// ignore empty string
		if b.Key == "" {
			continue
		}
		elements = append(elements, Element{
			Key:   b.Key,
			Count: b.Count,
		})
	}
	return
}

//Counter - the spacesaving counter
type Counter struct {
	list    []Element
	hash    map[string]uint32 //value:indexOfKeyInTheList
	isSuper bool
}

//Element - the bucket to hold each element
type Element struct {
	Key   string
	Count uint64
	//SubCounters e.g: Hashtag is the supercounter, it contains the keyword of hashtag itsself, as well as
	// subcounters of word and timezone associated with the key
	subCounters []*Counter
}
