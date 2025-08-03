package gocli

import (
	"github.com/dimonrus/gohelp"
	"testing"
)

func TestArguments_GetByName(t *testing.T) {
	t.Run("list", func(t *testing.T) {
		args := Arguments{{
			Type:  "int",
			Value: gohelp.Ptr(10),
			Label: "Number",
			Name:  "count",
		}, {
			Type:  "string",
			Value: gohelp.Ptr("some"),
			Label: "Name",
			Name:  "name",
		}, {
			Type:  "bool",
			Value: gohelp.Ptr(true),
			Label: "IsEnabled",
			Name:  "enabled",
		}, {
			Type:  "float",
			Value: gohelp.Ptr(33.44),
			Label: "Part",
			Name:  "part",
		}}
		arg := args.GetByName("enabled")
		if arg == nil || arg.GetBool() != true {
			t.Fatal("wrong logic for bool")
		}
		arg = args.GetByName("name")
		if arg == nil || arg.GetString() != "some" {
			t.Fatal("wrong logic for string")
		}
		arg = args.GetByName("part")
		if arg == nil || arg.GetFloat() != 33.44 {
			t.Fatal("wrong logic for float")
		}
	})
}

func TestArgumentMap_ToList(t *testing.T) {
	am := ArgumentMap{
		"count": {
			Type:  "int",
			Value: gohelp.Ptr(10),
			Label: "Number",
		},
		"name": {
			Type:  "string",
			Value: gohelp.Ptr("some"),
			Label: "Name",
		},
	}
	list := am.ToList()
	nameArgument := list.GetByName("name")
	if nameArgument == nil || nameArgument.Name != "name" {
		t.Fatal("wrong logic")
	}
}
