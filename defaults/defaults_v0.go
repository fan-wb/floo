package defaults

import (
	"errors"
	"net/url"
	"runtime"

	"github.com/sahib/config"
)

// DaemonDefaultURL returns the default URL for the current OS.
func DaemonDefaultURL() string {
	// If the platform supports unix sockets,
	// we should make use of it.
	switch runtime.GOOS {
	case "linux", "darwin":
		// See "man 7 unix" - we use an abstract unix domain socket.
		// This means there is no socket file on the file system.
		// (other tools use unix:@/path, but Go does not support that notation)
		// This also means that there are no user rights on the socket file.
		// If you need this, specify the url in the config.
		return "unix:/tmp/floo.socket?abstract=true"
	default:
		return "tcp://127.0.0.1:6666"
	}
}

func urlValidator(val interface{}) error {
	s, ok := val.(string)
	if !ok {
		return errors.New("url is not an string")
	}

	_, err := url.Parse(s)
	return err
}

var DefaultsV0 = config.DefaultMapping{
	"daemon": config.DefaultMapping{
		"url": config.DefaultEntry{
			Default:      DaemonDefaultURL(),
			NeedsRestart: true,
			Docs:         "URL of the daemon process.",
			Validator:    urlValidator,
		},
		"ipfs_path_or_url": config.DefaultEntry{
			Default:      "",
			NeedsRestart: true,
			Docs:         "URL or path to the IPFS repository you want to use.",
		},
		"enable_pprof": config.DefaultEntry{
			Default:      true,
			NeedsRestart: true,
			Docs:         "Enable a pprof profile server on startup (see < floo d p --help >)",
		},
	},
}
