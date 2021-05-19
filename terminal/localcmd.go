package terminal

import (
	"github.com/c-bata/go-prompt"
	"github.com/zhyhang/redis-client/platform"
	"os"
)

var localCmdFunMap = map[string]func(inputs *Inputs){
	"quit":  exitShell,
	"exit":  exitShell,
	"clear": clearShell,
	"cls":   clearShell,
}

func exitShell(inputs *Inputs) {
	err := tunnel.Destroy()
	platform.HandleExit()
	if err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}

func clearShell(inputs *Inputs) {
	writer := prompt.NewStdoutWriter()
	writer.EraseScreen()
	writer.CursorGoTo(0, 0)
	_ = writer.Flush()
}
