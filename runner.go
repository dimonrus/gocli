package gocli

import "github.com/dimonrus/porterr"

// Runner callback application starter
type Runner func(args ...Argument) (e porterr.IError)

// RunnerList list of application runners
type RunnerList []Runner
