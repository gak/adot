package main

import (
	"fmt"
	"github.com/alecthomas/kong"
	"github.com/gak/adot"
	"github.com/pkg/errors"
)

type Arg struct {
	Init struct{
		URL string `arg`
	} `cmd`
	Existing struct{
		URL string `arg`
	} `cmd`
	Add struct{
		Path string `arg`
	} `cmd`
	Remove struct{
		Path string `arg`
	} `cmd`
	Push struct{} `cmd`
	Pull struct{} `cmd`
}

func main() {
	arg := Arg{}
	ctx := kong.Parse(&arg)
	a := &adot.ADot{}

	if err := a.Prepare(); err != nil {
		report(err)
	}

	err := execute(a, ctx, &arg)
	if err != nil {
		report(err)
	}
}

func execute(a *adot.ADot, ctx *kong.Context, arg *Arg) error {
	cmd := ctx.Command()
	switch cmd {
	case "new <url>":
		return a.InitNew(arg.Init.URL)
	case "existing <url>":
		return a.InitExisting(arg.Init.URL)
	case "add":
		return a.Add(arg.Add.Path)
	case "rm":
		return a.Remove(arg.Remove.Path)
	case "push":
		return a.Push()
	case "pull":
		return a.Pull()
	default:
		return errors.New(fmt.Sprintf("unhandled command: %v", cmd))
	}
}

func report(err error) {
	fmt.Printf("%+v\n", err)
}
