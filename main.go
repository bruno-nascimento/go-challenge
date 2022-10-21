package main

import (
	"go-challenge/cmd"
	"go-challenge/internal/args"
)

func main() {
	if config, err := args.Parse(); err == nil {
		cmd.Run(config)
	}
}
