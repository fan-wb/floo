package client

type Whoami struct {
	CurrentUser string
	Owner       string
	Fingerprint string
	IsOnline    bool
}

type RemoteFolder struct {
	Folder           string `yaml:"Folder"`
	ReadOnly         bool   `yaml:"ReadOnly"`
	ConflictStrategy string `yaml:"ConflictStrategy"`
}

type Remote struct {
	Name             string         `yaml:"Name"`
	Fingerprint      string         `yaml:"Fingerprint"`
	Folders          []RemoteFolder `yaml:"Folders,flow"`
	AutoUpdate       bool           `yaml:"AutoUpdate"`
	ConflictStrategy string         `yaml:"ConflictStrategy"`
	AcceptPush       bool           `yaml:"AcceptPush"`
}

func (cl *Client) Whoami() (*Whoami, error) {

}
