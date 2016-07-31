package index

// Summary holds index data
type Summary struct {
	table map[uint32]*Entry
}

// Add and entry to index
func (s *Summary) Add(e *Entry) {
	s.table[e.Key] = e
}

// LookUp search in index table
func (s *Summary) LookUp(key uint32) (*Entry, bool) {
	value, ok := s.table[key]
	return value, ok
}

// GetTable returns the index holder
func (s *Summary) GetTable() map[uint32]*Entry {
	return s.table
}

// NewSummary returns new Summary
func NewSummary() *Summary {
	return &Summary{
		table: make(map[uint32]*Entry),
	}
}
