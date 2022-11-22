package cmd

import (
	"context"
	"floo/client"
	"floo/defaults"
	"fmt"
	"github.com/sahib/config"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
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
		daemonURL, _ := guessDaemonURL(ctx)

		if startNew {
			logVerbose(ctx, "Using url %s to check for running daemon.", daemonURL)
		} else {
			logVerbose(ctx, "Using url %s to connect to existing daemon.", daemonURL)
		}

		// Check if the daemon is running already:
		ctl, err := client.Dial(context.Background(), daemonURL)
		if err == nil {
			defer ctl.Close()
			return handler(ctx, ctl)
		}

		if !startNew {
			// Daemon was not running, and we may not start a new one.
			return ExitCode{DaemonNotResponding, "Daemon not running"}
		}

		// Start the server & pass the password:
		folder, err := guessRepoFolder(ctx)
		if err != nil {
			return ExitCode{
				BadArgs,
				fmt.Sprintf("could not guess folder: %v", err),
			}
		}

		logVerbose(ctx, "starting new daemon in background, on folder '%s'", folder)

		ctl, err = startDaemon(ctx, folder, daemonURL)
		if err != nil {
			return ExitCode{
				DaemonNotResponding,
				fmt.Sprintf("Unable to start daemon: %v", err),
			}
		}

		// Run the actual handler:
		defer ctl.Close()
		return handler(ctx, ctl)
	}
}

func getExecutablePath() (string, error) {
	// NOTE: This might not work on other platforms.
	//       In this case we fall back to LookPath().
	exePath, err := os.Readlink("/proc/self/exe")
	if err != nil {
		return exec.LookPath("floo")
	}

	return filepath.Clean(exePath), nil
}

func startDaemon(ctx *cli.Context, repoPath, daemonURL string) (*client.Client, error) {
	stat, err := os.Stat(repoPath)
	if err != nil {
		return nil, err
	}

	if !stat.IsDir() {
		return nil, fmt.Errorf("< %s > is not a directory", repoPath)
	}

	exePath, err := getExecutablePath()
	if err != nil {
		return nil, err
	}

	logVerbose(ctx, "using executable path: %s", exePath)

	logVerbose(
		ctx,
		"No Daemon running at %s. Starting daemon from binary: %s",
		daemonURL,
		exePath,
	)

	daemonArgs := []string{
		"--repo", repoPath,
		"--url", daemonURL,
		"daemon", "launch",
	}

	argString := fmt.Sprintf("'%s'", strings.Join(daemonArgs, "' '"))
	logVerbose(ctx, "Starting daemon as: %s %s", exePath, argString)
	proc := exec.Command(exePath, daemonArgs...) // #nosec
	proc.Env = append(proc.Env, fmt.Sprintf("PATH=%s", os.Getenv("PATH")))
	if err := proc.Start(); err != nil {
		log.Infof("Failed to start the daemon: %v", err)
		return nil, err
	}

	// This will likely suffice for most cases:
	time.Sleep(500 * time.Millisecond)

	warningPrinted := false
	for i := 0; i < 500; i++ {
		ctl, err := client.Dial(context.Background(), daemonURL)
		if err != nil {
			// Only print this warning once...
			if !warningPrinted && i >= 100 {
				log.Warnf("waiting a bit long for daemon to bootup...")
				warningPrinted = true
			}
			time.Sleep(50 * time.Millisecond)
			continue
		}
		return ctl, nil
	}
	return nil, fmt.Errorf("daemon could not be started or took to long")
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
	cfg, err := openConfig(folder)
	if err != nil {
		return ctx.GlobalString("url"), nil
	}

	return cfg.String("daemon.url"), nil
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

func openConfig(folder string) (*config.Config, error) {
	configPath := filepath.Join(folder, "config.yml")
	cfg, err := defaults.OpenMigratedConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("couldn't find a config: %v", err)
	}
	return cfg, nil
}
