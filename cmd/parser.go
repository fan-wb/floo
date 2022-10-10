package cmd

import (
	"floo/version"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"strings"
)

func RunCmdline(args []string) int {
	app := cli.NewApp()
	app.Name = "floo"
	app.Usage = "decentralized and secure file synchronization"
	app.EnableBashCompletion = true
	app.Version = fmt.Sprintf(
		"%s [buildtime: %s] (client version)",
		version.String(),
		version.BuildTime,
	)
	//app.CommandNotFound = commandNotFound
	app.Description = "floo can be used to securely store, version and synchronize files among peers."

	// Groups:
	repoGroup := formatGroup("repository")
	wdirGroup := formatGroup("working tree")
	vcscGroup := formatGroup("version control")
	netwGroup := formatGroup("network")

	app.BashComplete = func(ctx *cli.Context) {
		for _, cmd := range app.Commands {
			fmt.Println(cmd.Name)
		}
	}

	app.Commands = TranslateHelp([]cli.Command{
		{
			Name:     "init",
			Category: repoGroup,
			Action:   handleInit,
		}, {
			Name:     "whoami",
			Aliases:  []string{"id"},
			Category: netwGroup,
			// withDaemon
			Action: handleWhoami,
		}, {
			Name:     "remote",
			Aliases:  []string{"rmt", "r"},
			Category: netwGroup,
			Subcommands: []cli.Command{
				{
					Name:    "add",
					Aliases: []string{"a", "set"},
					Action:  withArgCheck(needAtLeast(2), withDaemon(handleRemoteAdd, true)),
				},
			},
		},
	})

	exitCode := Success
	if err := app.Run(args); err != nil {
		//log.Error(prettyPrintError(err))
		log.Error(err)
		cerr, ok := err.(ExitCode)
		if !ok {
			exitCode = UnknownError
		}

		exitCode = cerr.Code
	}

	return exitCode
}

func formatGroup(category string) string {
	return "\n" + strings.ToUpper(category) + " COMMANDS"
}
