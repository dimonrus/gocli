package gocli

const (
	ArgumentTypeString = "string"
	ArgumentTypeInt    = "int"
	ArgumentTypeUint   = "uint"
	ArgumentTypeBool   = "bool"
	ArgumentTypeFloat  = "float"
)

// ArgumentMap Console app arguments
type ArgumentMap map[string]Argument

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

// GetString Get string value of argument
func (a Argument) GetString() string {
	value := a.Value.(*string)
	return *value
}

// GetInt Get int value of argument
func (a Argument) GetInt() int64 {
	value := a.Value.(*int64)
	return *value
}

// GetUnit Get int value of argument
func (a Argument) GetUnit() uint64 {
	value := a.Value.(*uint64)
	return *value
}

// GetBool Get bool value of argument
func (a Argument) GetBool() bool {
	value := a.Value.(*bool)
	return *value
}

// GetFloat Get float value of argument
func (a Argument) GetFloat() float64 {
	value := a.Value.(*float64)
	return *value
}

// Arguments list of arguments
type Arguments []Argument

// GetByName get argument by name
func (a *Arguments) GetByName(name string) *Argument {
	if a == nil {
		return nil
	}
	for i := range *a {
		argument := (*a)[i]
		if argument.Name == name {
			return &argument
		}
	}
	return nil
}
