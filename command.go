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
	CommandPrefix   = "-"
	CommandAssignee = "="
)

// Command is an argument list
type Command struct {
	arguments  []Argument
	connection net.Conn
	origin     []byte
	m          sync.RWMutex
}

// Result of command to connection
func (c Command) Result(result []byte) porterr.IError {
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
func (c Command) Arguments() []Argument {
	c.m.RLock()
	defer c.m.RUnlock()
	return c.arguments
}

// Get origin command
func (c Command) GetOrigin() string {
	c.m.RLock()
	defer c.m.RUnlock()
	return string(c.origin)
}

// Render command
func (c Command) String() string {
	c.m.RLock()
	defer c.m.RUnlock()
	var command []string
	for _, a := range c.arguments {
		command = append(command, a.Name)
	}
	return strings.Join(command, " ")
}

// Parse command
func ParseCommand(command []byte) Command {
	sCommand := strings.Trim(string(command), "	 \n")
	cmd := Command{
		arguments: make([]Argument, 0),
		origin:    command,
	}
	var words []string
	var word string
	for _, c := range sCommand {
		if c == ' ' {
			if len(word) > 0 {
				words = append(words, strings.Split(word, CommandAssignee)...)
				word = ""
			}
		} else if c != '\n' && c != '\t' {
			word += string(c)
		}
	}
	if len(word) > 0 {
		words = append(words, strings.Split(word, CommandAssignee)...)
		word = ""
	}
	for i := range words {
		words[i] = strings.Trim(words[i], CommandPrefix)
		var argument Argument
		argument.Name = words[i]

		if v, err := strconv.ParseInt(words[i], 10, 64); err == nil {
			argument.Type = ArgumentTypeInt
			argument.Value = &v
		} else if v, err := strconv.ParseBool(words[i]); err == nil {
			argument.Type = ArgumentTypeBool
			argument.Value = &v
		} else if v, err := strconv.ParseUint(words[i], 10, 64); err == nil {
			argument.Type = ArgumentTypeUint
			argument.Value = &v
		} else {
			argument.Type = ArgumentTypeString
			argument.Value = &words[i]
		}

		cmd.arguments = append(cmd.arguments, argument)
	}
	return cmd
}
