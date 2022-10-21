package cmd

import (
	"go-challenge/cmd/generate"
	"go-challenge/cmd/summary"
	"go-challenge/internal/args"
)

func Run(config *args.Config) {
	if config.GenerateTestFiles {
		generate.Cmd(config)
		return
	}
	summary.Cmd(config)
}
