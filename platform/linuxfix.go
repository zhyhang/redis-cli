// +build !windows

package platform

import (
	"os"
	"os/exec"
)

func HandleExit() {
	// workaround for the issue: Terminal does not display command after go-prompt exit
	// https://github.com/c-bata/go-prompt/issues/228
	rawModeOff := exec.Command("/bin/stty", "-raw", "echo")
	rawModeOff.Stdin = os.Stdin
	_ = rawModeOff.Run()
	rawModeOff.Wait()
}
