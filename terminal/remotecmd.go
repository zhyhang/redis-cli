package terminal

func updateConf(input *ShellInputs, config *shellConfig) {
	if "monitor" == input.Cmd {
		config.modeMonitor = true
	}
}
