package cmd

import (
	"fmt"
	"github.com/urfave/cli"
	"strings"
)

type helpEntry struct {
	Usage       string
	ArgsUsage   string
	Description string
	Complete    cli.BashCompleteFunc
	Flags       []cli.Flag
}

var helpTexts = map[string]helpEntry{
	"init": {
		Usage:     "Initialize a new repository.",
		ArgsUsage: "<username>",
		Complete:  completeArgsUsage,
	},
}

func translateHelp(cmds []cli.Command, prefix []string) {
	for idx := range cmds {
		path := append(append([]string{}, prefix...), cmds[idx].Name)
		injectHelp(&cmds[idx], strings.Join(path, "."))
		translateHelp(cmds[idx].Subcommands, path)
	}
}

func injectHelp(cmd *cli.Command, path string) {
	help, ok := helpTexts[path]
	if !ok {
		panic(fmt.Sprintf("bug: no such help entry: %v", path))
	}

	cmd.Usage = help.Usage
	cmd.ArgsUsage = help.ArgsUsage
	cmd.Description = help.Description
	cmd.BashComplete = help.Complete
	cmd.Flags = help.Flags
}

// TranslateHelp fills in the usage and description for each command.
// This is separated from the command definition to make things more readable,
// and separate logic from the (lengthy) documentation.
func TranslateHelp(cmds []cli.Command) []cli.Command {
	translateHelp(cmds, nil)
	return cmds
}
