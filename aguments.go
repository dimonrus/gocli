package gocli

const (
	ArgumentTypeString  = "string"
	ArgumentTypeInt     = "int"
	ArgumentTypeUint    = "uint"
	ArgumentTypeBool    = "bool"
)

// Console app arguments
type Arguments map[string]Argument

// Argument
type Argument struct {
	Type  string
	Value interface{}
	Label string
}

// Get string value of argument
func (a Argument) GetString() string {
	value := a.Value.(*string)
	return *value
}

// Get int value of argument
func (a Argument) GetInt() int64 {
	value := a.Value.(*int64)
	return *value
}

// Get int value of argument
func (a Argument) GetUnit() uint64 {
	value := a.Value.(*uint64)
	return *value
}

// Get bool value of argument
func (a Argument) GetBool() bool {
	value := a.Value.(*bool)
	return *value
}
