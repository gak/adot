package main

import (
	"fmt"
	"github.com/alecthomas/kong"
	"github.com/gak/adot"
	"github.com/pkg/errors"
)

type Arg struct {
	Init string `cmd`
	Push string `cmd`
	Pull string `cmd`
}

func main() {
	arg := Arg{}
	ctx := kong.Parse(&arg)
	a := &adot.ADot{}

	if err := a.Defaults(); err != nil {
		panic(err)
	}

	err := execute(a, ctx)
	if err != nil {
		panic(err)
	}
}
func execute(a *adot.ADot, ctx *kong.Context) error {
	cmd := ctx.Command()
	switch cmd {
	case "init":
		return a.Init()
	case "push":
		return a.Push()
	case "pull":
		return a.Pull()
	default:
		return errors.New(fmt.Sprintf("unhandled command: %v", cmd))
	}
}
