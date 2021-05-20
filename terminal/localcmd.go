package terminal

import (
	"fmt"
	"github.com/c-bata/go-prompt"
	"github.com/zhyhang/redis-client/platform"
	"github.com/zhyhang/redis-client/redis"
	"github.com/zhyhang/redis-client/util"
	"os"
	"strconv"
	"strings"
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
	continueRemote = false
	if len(inputs.Args) == 0 {
		fmt.Println(util.ShellHelp)
		return
	}
	arg0 := strings.ToLower(inputs.Args[0])
	if arg0 == util.ShellHelpAll {
		for i := 0; i < len(util.GetCommandHelps()); i++ {
			printGroupCmdHelp(i)
		}
		return
	}
	if arg0[0] == '@' {
		findPrintGroupCmdHelp(arg0[1:])
		return
	}
	return
}

func findPrintGroupCmdHelp(group string) {
	groups := util.GetCommandGroups()
	for i, g := range groups {
		if group == g {
			printGroupCmdHelp(i)
			return
		}
	}
}

func printGroupCmdHelp(gi int) {
	helps := util.GetCommandHelps()
	for _, h := range helps {
		if h.Group == gi {
			printCmdHelp(h)
		}
	}
}

func printCmdHelp(cmdHelp util.CommandHelp) {
	fmt.Printf("\n  %s %s\n  summary: %s\n  since: %s\n", cmdHelp.Name, cmdHelp.Params, cmdHelp.Summary, cmdHelp.Since)
}
