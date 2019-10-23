package spacesaving

import (
	"log"
	"sort"
)

//NewMHCounter - initialise memory hungry counter
func NewMHCounter(s int, isSuperCounter bool) *MHCounter {
	mh := MHCounter{
		buckets: make(map[string]*Element, 0),
		isSuper: isSuperCounter,
	}
	return &mh
}

//Hit - for every instance in stream, update the counter
func (mh *MHCounter) Hit(key string) {

	if elem, found := mh.buckets[key]; found {
		elem.Count = elem.Count + 1
	} else {
		bucket := &Element{
			Count: 1,
		}
		// only create subcounters for supercounter
		if mh.isSuper {
			bucket.subCounters = make([]Counter, numSubCounters)
			for i := 0; i < numSubCounters; i++ {
				bucket.subCounters[i] = NewMHCounter(1000, false)
			}
		}
		mh.buckets[key] = bucket
	}
}

//GetAll return all elements in the counter
func (mh *MHCounter) GetAll() (elements ElementList) {
	elements = make(ElementList, len(mh.buckets))
	// output from higt to low count
	i := 0
	for k, v := range mh.buckets {
		elements[i] = Element{k, v.Count, nil}
		i++
	}
	sort.Sort(sort.Reverse(elements))
	return elements
}

//GetSubCounter - get subcounter of the bucket with key
func (mh *MHCounter) GetSubCounter(key string, i int) (subCounter Counter) {
	if i >= numSubCounters {
		log.Panicln("subcounter index out of bound")
	}
	if elem, found := mh.buckets[key]; found {
		return elem.subCounters[i]
	}
	log.Panicln("call for subcounter under no-exist supercounter's bucket")
	return
}

//MHCounter - the memory hungry counter
type MHCounter struct {
	buckets map[string]*Element
	isSuper bool
}
