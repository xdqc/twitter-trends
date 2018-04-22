package spacesaving

//NewCounter - initialise a counter
func NewCounter(size int, isSuper bool) Counter {
	if size == 0 {
		return NewMHCounter(size, isSuper)
	}
	return NewSSCounter(size, isSuper)
}

//Counter - general counter
type Counter interface {
	Hit(string)
	GetAll() ElementList
	GetSubCounter(string, int) Counter
}

//Element - the bucket to hold each element
type Element struct {
	Key   string
	Count uint64
	//SubCounters e.g: Hashtag is the supercounter, it contains the keyword of hashtag itsself, as well as
	// subcounters of word and timezone associated with the key
	subCounters []Counter
}

//ElementList - list of element, sortable
type ElementList []Element

func (es ElementList) Len() int           { return len(es) }
func (es ElementList) Less(i, j int) bool { return es[i].Count < es[j].Count }
func (es ElementList) Swap(i, j int)      { es[i], es[j] = es[j], es[i] }
