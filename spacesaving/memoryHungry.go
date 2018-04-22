package spacesaving

import (
	"sort"
)

//NewMHCounter - initialise memory hungry counter
func NewMHCounter(s int, isSuperCounter bool) *MHCounter {
	mh := MHCounter{
		buckets: make(map[string]uint64, 0),
	}
	return &mh
}

//Hit - for every instance in stream, update the counter
func (mh *MHCounter) Hit(key string) {
	if count, found := mh.buckets[key]; found {
		mh.buckets[key] = count + 1
	} else {
		mh.buckets[key] = 0
	}
}

//GetAll return all elements in the counter
func (mh *MHCounter) GetAll() (elements ElementList) {
	elements = make(ElementList, len(mh.buckets))
	// output from higt to low count
	i := 0
	for k, v := range mh.buckets {
		elements[i] = Element{k, v, nil}
		i++
	}
	sort.Sort(sort.Reverse(elements))
	return elements
}

//GetSubCounter - get subcounter of the bucket with key
func (mh *MHCounter) GetSubCounter(key string, i int) (subCounter *SSCounter) {
	return nil
}

//MHCounter - the memory hungry counter
type MHCounter struct {
	buckets map[string]uint64
}

//ElementList - list of element
type ElementList []Element

func (es ElementList) Len() int           { return len(es) }
func (es ElementList) Less(i, j int) bool { return es[i].Count < es[j].Count }
func (es ElementList) Swap(i, j int)      { es[i], es[j] = es[j], es[i] }
