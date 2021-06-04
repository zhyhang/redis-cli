package terminal

type CmdFlags struct {
	Host        string
	Port        int
	ReturnError bool
}

type ShellInputs struct {
	LineTrim string `json:"line_trim"`
	RawCmd   string `json:"raw_cmd"`
	// lower case first one word of input
	Cmd  string   `json:"cmd"`
	Args []string `json:"args"`
}

type shellConfig struct {
	cmdLine     *CmdFlags
	modeMonitor bool
}

func NewCmdFlags() *CmdFlags {
	return &CmdFlags{Host: "127.0.0.1", Port: 6379, ReturnError: false}
}

func newShellConfig() *shellConfig {
	return &shellConfig{}
}
