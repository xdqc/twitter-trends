package spacesaving

import (
	"log"
)

var (
	size int
	//`2` indicates the two sub-counters for timezone and word associated with the hashtag, respectively
	numSubCounters = 2
	// only create subcounters for top N heavy hitters (save memory)
	numTopHeavyHitterThatHasSubcounter = 1000
)

//NewSSCounter - initialise spacesaving counter
func NewSSCounter(s int, isSuperCounter bool) *SSCounter {
	if isSuperCounter {
		size = s
	}
	ss := SSCounter{
		list:    make([]*Element, s),
		hash:    make(map[string]uint32, s),
		isSuper: isSuperCounter,
	}

	for i := 0; i < s; i++ {
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

		// only create subcounters for top N heavy hitters
		if ss.isSuper && bucket.subCounters == nil && int(idx) >= len(ss.list)-numTopHeavyHitterThatHasSubcounter {
			log.Println("Create subctrs for top#", len(ss.list)-int(idx), "heavy hitter", key)
			bucket.subCounters = make([]Counter, numSubCounters)
			if bucket.Key == " " {
				// Use same size subcounter as supercouter to count tweets with no hashtag
				bucket.subCounters[0] = NewSSCounter(1000, false)    //timezone
				bucket.subCounters[1] = NewSSCounter(size*10, false) //word
				log.Println("subcounter created for non-hashtag tweet.")
			} else {
				for i := 0; i < numSubCounters; i++ {
					// Use numTopHeavyHitterThatHasSubcounter sized subcounter
					bucket.subCounters[i] = NewSSCounter(numTopHeavyHitterThatHasSubcounter, false)
				}
			}
		}
	} else {
		// new element => replace the first element(lowest count) with new key
		idx = 0
		bucket = ss.list[idx]
		// create subcounters to replace the old ones in buckets whose key been overwritten
		if ss.isSuper && bucket.subCounters != nil {
			// counter 'full' indicator
			log.Println("create subcounters for new", key, "to replace", bucket.Key, "with count", bucket.Count)
			bucket.subCounters = make([]Counter, numSubCounters)
			for i := 0; i < numSubCounters; i++ {
				bucket.subCounters[i] = NewSSCounter(1000, false)
			}
		}
		delete(ss.hash, bucket.Key)
		bucket.Key = key
		ss.hash[key] = idx
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
		if len(elem.subCounters) > 0 {
			return elem.subCounters[i]
		}
		return nil
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
