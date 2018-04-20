package spacesaving

import (
	"log"
	"sync"
)

//Counter - the spacesaving counter
type Counter struct {
	list []Element
	hash map[string]uint32
}

//Element - the buckt to hold each element
type Element struct {
	Key   string
	Count uint64
	//SubCounters e.g: Hashtag contains the keyword of hashtag itsself, as well as
	// sub-counters of word and timezone associated with the key
	subCounters []*Counter
}

var (
	//use 1000 or 10000 for the number of element buckets of the counter
	size int
	//`2` indicates the two sub-counters for timezone and word associated with the hashtag, respectively
	numSubCounters = 2
	//Use mutex solve concurrent map read and map write problem
	mutex = sync.RWMutex{}
)

//NewCounter initialise spacesaving counter
func NewCounter(s int) *Counter {
	size = s
	ss := Counter{
		list: make([]Element, size),
		hash: make(map[string]uint32, size),
	}
	return &ss
}

//Hit for every instance in stream
func (ss *Counter) Hit(key string) {
	var (
		idx    uint32
		found  bool
		bucket *Element
	)

	if idx, found = ss.hash[key]; found {
		bucket = &ss.list[idx]
	} else {
		// replace the lowest count with new element
		idx = 0
		bucket = &ss.list[idx]
		delete(ss.hash, bucket.Key)

		ss.hash[key] = idx
		bucket.Key = key
		bucket.subCounters = make([]*Counter, 2)
		for i := 0; i < numSubCounters; i++ {
			bucket.subCounters[i] = NewCounter(size)
		}
	}

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

//GetSubCounter get i-th subcounter of the bucket with key
func (ss *Counter) GetSubCounter(key string, i int) (subCounter *Counter) {
	if i >= numSubCounters {
		log.Panicln("sub-counter index out of bound")
	}
	// mutex.Lock()
	if elemIdx, found := ss.hash[key]; found {
		elem := ss.list[elemIdx]
		return elem.subCounters[i]
	}
	// mutex.Unlock()
	log.Panicln("call for subcounter of no-exist supercounter")
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

//Reset reset the spacesaving counter
func (ss *Counter) Reset() {
	empty := Element{}
	for i := range ss.list {
		delete(ss.hash, ss.list[i].Key)
		ss.list[i] = empty
	}
}
