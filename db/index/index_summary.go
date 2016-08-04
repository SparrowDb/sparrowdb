package index

// Summary holds index data
type Summary struct {
	table map[uint32]*Entry
	count uint32
}

// Add and entry to index
func (s *Summary) Add(e *Entry) {
	s.count++
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

// Count returns the number of itens in Summary
func (s *Summary) Count() uint32 {
	return s.count
}

// NewSummary returns new Summary
func NewSummary() *Summary {
	return &Summary{
		table: make(map[uint32]*Entry),
	}
}
