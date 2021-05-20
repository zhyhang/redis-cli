package terminal

import (
	"fmt"
	"github.com/c-bata/go-prompt"
	"github.com/zhyhang/redis-client/platform"
	"github.com/zhyhang/redis-client/redis"
	"github.com/zhyhang/redis-client/util"
	"os"
	"strconv"
)

var localCmdFunMap = map[string]func(inputs *ShellInputs) (continueRemote bool){
	"quit":    exitShell,
	"exit":    exitShell,
	"clear":   clearShell,
	"cls":     clearShell,
	"connect": connectRedis,
	"help":    shellHelp,
	"?":       shellHelp,
}

func exitShell(inputs *ShellInputs) (continueRemote bool) {
	tunnel.Destroy()
	platform.HandleExit()
	continueRemote = false
	os.Exit(0)
	return
}

func clearShell(inputs *ShellInputs) (continueRemote bool) {
	writer := prompt.NewStdoutWriter()
	writer.EraseScreen()
	writer.CursorGoTo(0, 0)
	_ = writer.Flush()
	return false
}

func connectRedis(inputs *ShellInputs) (continueRemote bool) {
	iLen := len(inputs.Args)
	if iLen > 0 && inputs.Args[0] != "" {
		continueRemote = false
		port := 6379
		if iLen > 1 && inputs.Args[1] != "" {
			p, err := strconv.Atoi(inputs.Args[1])
			if err != nil {
				addr := inputs.Args[0] + ":" + inputs.Args[1]
				fmt.Println(redis.NotLinkMsg(addr))
				return
			}
			port = p
		}
		tun := redis.Establish(inputs.Args[0], port)
		if tun.Linked {
			tunnel = tun
		}
	} else {
		continueRemote = true
	}
	return
}

func shellHelp(inputs *ShellInputs) (continueRemote bool) {
	fmt.Println(util.ShellHelp)
	return false
}
