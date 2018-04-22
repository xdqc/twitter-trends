package spacesaving

//Counter -
type Counter interface {
	Hit(string)
	GetAll() ElementList
	GetSubCounter(string, int) *SSCounter
}

//NewCounter - initialise a counter
func NewCounter(size int, isSuper bool) Counter {
	if size == 0 {
		return NewMHCounter(size, isSuper)
	}
	return NewSSCounter(size, isSuper)
}
