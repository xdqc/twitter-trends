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

//NewSSCounter - initialise spacesaving counter
func NewSSCounter(s int, isSuperCounter bool) *SSCounter {
	size = s
	ss := SSCounter{
		list:    make([]*Element, size),
		hash:    make(map[string]uint32, size),
		isSuper: isSuperCounter,
	}

	for i := 0; i < size; i++ {
		ss.list[i] = &Element{}
	}

	return &ss
}

//Hit - for every instance in stream, update the counter
func (ss *SSCounter) Hit(key string) {
	var (
		idx    uint32
		found  bool
		bucket *Element
	)

	if idx, found = ss.hash[key]; found {
		// exsisting element => increment count
		bucket = ss.list[idx]
	} else {
		// new element => replace the first element(lowest count) with new key
		idx = 0
		bucket = ss.list[idx]
		delete(ss.hash, bucket.Key)
		bucket.Key = key
		ss.hash[key] = idx
		// only create subcounters for supercounter
		if ss.isSuper {
			bucket.subCounters = make([]Counter, numSubCounters)
			for i := 0; i < numSubCounters; i++ {
				bucket.subCounters[i] = NewSSCounter(size, false)
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

		b1 := ss.list[idx]
		b2 := ss.list[idx+1]

		//ignore counting ties
		if b1.Count <= b2.Count {
			break
		}

		//switch buckets pointer
		ss.hash[b1.Key] = idx + 1
		ss.hash[b2.Key] = idx
		ss.list[idx], ss.list[idx+1] = b2, b1
		idx++
	}
}

//GetSubCounter - get subcounter of the bucket with key
func (ss *SSCounter) GetSubCounter(key string, i int) (subCounter Counter) {
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
func (ss *SSCounter) GetAll() (elements ElementList) {
	elements = make([]Element, 0, len(ss.hash))
	// output from higt to low count
	for i := len(ss.list) - 1; i >= 0; i-- {
		b := ss.list[i]
		// ignore empty
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

//SSCounter - the spacesaving counter
type SSCounter struct {
	list    []*Element
	hash    map[string]uint32 //value:indexOfKeyInTheList
	isSuper bool
}
