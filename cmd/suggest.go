package cmd

import (
	"fmt"
	"github.com/urfave/cli"
	"strings"
)

func completeArgsUsage(ctx *cli.Context) {
	if command := findCurrentCommand(ctx); command != nil {
		if len(command.Flags) == 0 {
			return
		}

		for _, flag := range command.Flags {
			split := strings.SplitN(flag.GetName(), ",", 2)
			longName := split[0]
			fmt.Printf("--%s\n", longName)
		}

		fmt.Println(command.ArgsUsage)
	}

}

func findCurrentCommand(ctx *cli.Context) *cli.Command {
	for {
		par := ctx.Parent()
		if par == nil {
			break
		}
		ctx = par
	}
	var command *cli.Command
	for args := ctx.Args(); len(args) > 0; {
		subCommand := ctx.App.Command(args[0])
		args = args[1:]
		if subCommand != nil {
			command = subCommand
		}
	}
	return command
}
