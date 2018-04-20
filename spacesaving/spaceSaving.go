package spacesaving

import (
	"log"
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
	// e.g: Hashtag contains the keyword of hashtag itsself, as well as
	// sub-counters of word and timezone associated with the key
	SubCounters []*Counter
}

var (
	numSubCounters int
	size           int
)

//NewCounter initialise spacesaving counter
func NewCounter(s int, n int) *Counter {
	size = s
	numSubCounters = n
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
		bucket.SubCounters = make([]*Counter, 2)
		for i := 0; i < numSubCounters; i++ {
			bucket.SubCounters[i] = NewCounter(size, 2)
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
	if elemIdx, found := ss.hash[key]; found {
		elem := ss.list[elemIdx]
		return elem.SubCounters[i]
	}
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
