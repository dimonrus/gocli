package gocli

import (
	"fmt"
	"github.com/dimonrus/gohelp"
	"github.com/dimonrus/porterr"
	"net"
	"strconv"
	"strings"
	"sync"
)

const (
	CommandPrefix    = "-"
	CommandAssignee  = "="
	CommandDelimiter = ";"
)

var (
	ignored  = []rune{' ', '\n', '\t', '\r', '-'}
	assigned = '='
)

// Command is an argument list
type Command struct {
	// list of arguments
	arguments []Argument
	// net connection
	connection net.Conn
	// original command
	origin []byte
	// mutex for async access
	m sync.RWMutex
}

// Result of command to connection
func (c *Command) Result(result []byte) porterr.IError {
	c.m.Lock()
	defer c.m.Unlock()
	if c.connection == nil {
		return nil
	}
	data := fmt.Sprintf(gohelp.AnsiBlue+"--->: "+gohelp.AnsiGreen+"%s"+gohelp.AnsiReset, result)
	_, err := c.connection.Write([]byte(data))
	if err != nil {
		return porterr.New(porterr.PortErrorIO, "Result command write error: "+err.Error())
	}
	return nil
}

// Bind connection to command
func (c *Command) BindConnection(conn net.Conn) {
	c.m.Lock()
	defer c.m.Unlock()
	c.connection = conn
}

// UnBind connection
func (c *Command) UnbindConnection() {
	c.m.Lock()
	defer c.m.Unlock()
	c.connection = nil
}

// UnBind connection
func (c *Command) Arguments() []Argument {
	c.m.RLock()
	defer c.m.RUnlock()
	return c.arguments
}

// Get origin command
func (c *Command) GetOrigin() string {
	c.m.RLock()
	defer c.m.RUnlock()
	return string(c.origin)
}

// Render command
func (c *Command) String() string {
	c.m.RLock()
	defer c.m.RUnlock()
	var command []string
	for _, a := range c.arguments {
		command = append(command, a.Name)
	}
	return strings.Join(command, " ")
}

// Parse command
func ParseCommand(command []byte) *Command {
	var isIgnored, isAssignee bool
	var word []byte
	var argument Argument
	var valueInt64 int64
	var valueBool bool
	var valueUint64 uint64
	var err error
	var cmd = Command{
		arguments: make([]Argument, 0, 16),
		origin:    command,
	}
	var l = len(command) - 1
	for j, c := range command {
		isIgnored, isAssignee = false, false
		for i := 0; i < len(ignored); i++ {
			if ignored[i] == rune(c) {
				isIgnored = true
				break
			}
			if assigned == rune(c) {
				isAssignee = true
			}
		}
		if !isIgnored && !isAssignee {
			word = append(word, c)
			if j != l {
				continue
			}
		} else {
			if len(word) == 0 {
				continue
			}
		}
		argument.Name = string(word)
		if valueInt64, err = strconv.ParseInt(argument.Name, 10, 64); err == nil {
			argument.Type = ArgumentTypeInt
			argument.Value = &valueInt64
		} else if valueBool, err = strconv.ParseBool(argument.Name); err == nil {
			argument.Type = ArgumentTypeBool
			argument.Value = &valueBool
		} else if valueUint64, err = strconv.ParseUint(argument.Name, 10, 64); err == nil {
			argument.Type = ArgumentTypeUint
			argument.Value = &valueUint64
		} else {
			argument.Type = ArgumentTypeString
			argument.Value = &argument.Name
		}
		cmd.arguments = append(cmd.arguments, argument)
		word = nil
	}
	return &cmd
}
