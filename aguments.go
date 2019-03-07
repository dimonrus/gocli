package gocli

const (
	ArgumentTypeString = "string"
	ArgumentTypeInt    = "int"
)

// Console app arguments
type Arguments map[string]Argument

// Argument
type Argument struct {
	Type  string
	Value interface{}
	Label string
}

// Get string
func (a Argument) GetString() string {
	value := a.Value.(*string)
	return *value
}

// Get int value
func (a Argument) GetInt() int64 {
	value := a.Value.(*int64)
	return *value
}
