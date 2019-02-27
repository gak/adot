package main

import (
	"fmt"
	"github.com/alecthomas/kong"
	"github.com/gak/adot"
	"github.com/pkg/errors"
)

type Arg struct {
	New struct {
		URL string `arg`
	} `cmd help:"Create a new adot repository and push it."`
	Existing struct {
		URL string `arg`
	} `cmd help:"Clone an adot repository and load in the files."`
	Add struct {
		Path string `arg`
	} `cmd help:"Add a file to be tracked by adot."`
	Rm struct {
		Path string `arg`
	} `cmd help:"Remove a file from the adot repository. This will not remove the file from your home directory."`
	Push struct{} `cmd help:"Commit and push any changed files from your home directory."`
	Pull struct{} `cmd help:"Pull the latest repository and load in all the files."`
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
		return a.InitNew(arg.New.URL)
	case "existing <url>":
		return a.InitExisting(arg.Existing.URL)
	case "add":
		return a.Add(arg.Add.Path)
	case "rm":
		return a.Remove(arg.Rm.Path)
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
