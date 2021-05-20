package terminal

type CmdFlags struct {
	Host        string
	Port        int
	ReturnError bool
}

type ShellInputs struct {
	LineTrim string   `json:"line_trim"`
	RawCmd   string   `json:"raw_cmd"`
	Cmd      string   `json:"cmd"`
	Args     []string `json:"args"`
}

func NewCmdFlags() *CmdFlags {
	return &CmdFlags{Host: "127.0.0.1", Port: 6379, ReturnError: false}
}
