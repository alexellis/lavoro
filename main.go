package main

import (
	"github.com/alexellis/lavoro/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}
