package cmd

import (
	"floo/client"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"os"
	"path/filepath"
)

type ExitCode struct {
	Code    int
	Message string
}

func (err ExitCode) Error() string {
	return err.Message
}

type checkFunc func(ctx *cli.Context) int

func withArgCheck(checker checkFunc, handler cli.ActionFunc) cli.ActionFunc {
	return func(ctx *cli.Context) error {
		if checker(ctx) != Success {
			os.Exit(BadArgs)
		}
		return handler(ctx)
	}
}

func needAtLeast(min int) checkFunc {
	return func(ctx *cli.Context) int {
		if ctx.NArg() < min {
			log.Warningf("Need at least %d arguments.", min)
			if err := cli.ShowCommandHelp(ctx, ctx.Command.Name); err != nil {
				log.Warningf("Failed to display --help %v: ", err)
			}
			return BadArgs
		}
		return Success
	}
}

type cmdHandlerWithClient func(ctx *cli.Context, ctl *client.Client) error

func withDaemon(handler cmdHandlerWithClient, startNew bool) cli.ActionFunc {
	return func(ctx *cli.Context) error {

	}
}

func guessDaemonURL(ctx *cli.Context) (string, error) {
	if ctx.GlobalIsSet("url") {
		return ctx.GlobalString("url"), nil
	}

	folder, err := guessRepoFolder(ctx)
	if err != nil {
		log.Warnf("note: I don't know where the repository is or cannot read it.")
		log.Warnf("      I will continue with default values, cross fingers.")
		log.Warnf("      We recommend to set FLOO_PATH or pass --repo always.")
		log.Warnf("      Alternatively you can cd to your repository.")
		return ctx.GlobalString("url"), err
	}
}

func guessRepoFolder(ctx *cli.Context) (string, error) {
	if ctx.GlobalIsSet("repo") {
		return ctx.GlobalString("repo"), nil
	}

	guessLocations := []string{
		// TODO: just one now
		".",
	}
	var lastError error
	for _, guessLocation := range guessLocations {
		repoFolder := mustAbsPath(guessLocation)
		if _, err := os.Stat(filepath.Join(repoFolder, "config.yml")); err != nil {
			lastError = err
			continue
		}
		return repoFolder, nil
	}
	return "", lastError
}

func mustAbsPath(path string) string {
	absPath, err := filepath.Abs(path)
	if err != nil {
		fmt.Printf("Failed to get absolute repo path: %v", err)
		os.Exit(1)
	}
	return absPath
}

func openConfig(folder string) (*config.Config) {

}