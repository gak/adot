package main

import (
	"github.com/gak/adot"
)

func main() {
	a := &adot.ADot{}

	if err := a.Run(); err != nil {
		panic(err)
	}
}
