package spql

// Query holds query parsed content
type Query struct {
	Action  string
	Method  string
	Params  interface{}
	Filters map[string]string
}
