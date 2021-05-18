package terminal

type CmdFlags struct {
	Host string
	Port int
	ReturnError bool
}

func NewCmdFlags() *CmdFlags {
	return &CmdFlags{Host: "127.0.0.1",Port: 6379,ReturnError: false}
}
