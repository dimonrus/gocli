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
	ignored  = []rune{' ', '\n', '\t', '\r'}
	assigned = '='
	dash     = '-'
)

// Command is an argument list
type Command struct {
	// list of arguments
	arguments Arguments
	// net connection
	connection net.Conn
	// original command
	origin []byte
	// mutex for async access
	m sync.RWMutex
}

// Result of command to connection
func (c *Command) Result(result []byte) porterr.IError {
	data := fmt.Sprintf(gohelp.AnsiBlue+"--->: "+gohelp.AnsiGreen+"%s"+gohelp.AnsiReset, result)
	return c.Response([]byte(data))
}

// Response Flat result of command to connection
func (c *Command) Response(result []byte) porterr.IError {
	c.m.Lock()
	defer c.m.Unlock()
	if c.connection == nil {
		return nil
	}
	_, err := c.connection.Write(result)
	if err != nil {
		return porterr.New(porterr.PortErrorIO, "Result command write error: "+err.Error())
	}
	return nil
}

// BindConnection Bind connection to command
func (c *Command) BindConnection(conn net.Conn) {
	c.m.Lock()
	defer c.m.Unlock()
	c.connection = conn
}

// UnbindConnection UnBind connection
func (c *Command) UnbindConnection() {
	c.m.Lock()
	defer c.m.Unlock()
	c.connection = nil
}

// Arguments command arguments
func (c *Command) Arguments() Arguments {
	c.m.RLock()
	defer c.m.RUnlock()
	return c.arguments
}

// GetOrigin Get origin command
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

// ParseCommand Parse command
func ParseCommand(command []byte) *Command {
	var isIgnored, isAssignee bool
	var word = make([]byte, 0, 16)
	var k int
	cmd := Command{
		origin:    command,
		arguments: make(Arguments, 8),
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
		if len(word) == 0 && rune(c) == dash {
			isIgnored = true
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
		if k >= cap(cmd.arguments) {
			cmd.arguments = append(cmd.arguments, make(Arguments, 8)...)
		}
		isUint, isInt, isFloat, isBool, isString := gohelp.CheckTypeOf(word)
		value := string(word)
		switch true {
		case isUint:
			if valueUint64, err := strconv.ParseUint(value, 10, 64); err == nil {
				cmd.arguments[k].Type = ArgumentTypeUint
				cmd.arguments[k].Value = &valueUint64
			}
		case isInt:
			if valueInt64, err := strconv.ParseInt(value, 10, 64); err == nil {
				cmd.arguments[k].Type = ArgumentTypeInt
				cmd.arguments[k].Value = &valueInt64
			}
		case isFloat:
			if valueFloat64, err := strconv.ParseFloat(value, 64); err == nil {
				cmd.arguments[k].Type = ArgumentTypeFloat
				cmd.arguments[k].Value = &valueFloat64
			}
		case isBool:
			if valueBool, err := strconv.ParseBool(value); err == nil {
				cmd.arguments[k].Type = ArgumentTypeBool
				cmd.arguments[k].Value = &valueBool
			}
		case isString:
			cmd.arguments[k].Type = ArgumentTypeString
			cmd.arguments[k].Value = &value
		}
		cmd.arguments[k].Name = value
		k++
		word = word[:0]
	}
	cmd.arguments = cmd.arguments[:k]
	return &cmd
}
