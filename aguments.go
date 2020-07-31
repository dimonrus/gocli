package gocli

const (
	ArgumentTypeString = "string"
	ArgumentTypeInt    = "int"
	ArgumentTypeUint   = "uint"
	ArgumentTypeBool   = "bool"
)

// Console app arguments
type Arguments map[string]Argument

// Get command from arguments
func (a Arguments) Command() *Command {
	cmd := &Command{}
	for _, arg := range a {
		if arg.Value != nil {
			cmd.arguments = append(cmd.arguments, arg)
		}
	}
	return cmd
}

// Argument struct
type Argument struct {
	// Type of argument
	Type string
	// Value of argument
	Value interface{}
	// Label of argument
	Label string
	// Name of argument
	Name string
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
