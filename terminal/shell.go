package terminal

import (
	"fmt"
	prompt "github.com/c-bata/go-prompt"
	"github.com/zhyhang/redis-client/redis"
	"github.com/zhyhang/redis-client/util"
	"strings"
)

var tunnel *redis.Tunnel

func Run(flags *CmdFlags) {
	tunnel = redis.Establish(flags.Host, flags.Port)
	prefix, _ := changeLivePrefix()
	p := prompt.New(
		exec,
		suggest,
		prompt.OptionPrefix(prefix),
		prompt.OptionLivePrefix(changeLivePrefix),
		prompt.OptionTitle(util.ShellTitle),
		prompt.OptionMaxSuggestion(5),
		prompt.OptionPrefixTextColor(prompt.Green),
	)
	p.Run()
}

func exec(input string) {
	trimInput := strings.TrimSpace(input)
	if trimInput == "" {
		return
	}
	inputs := getInputs(trimInput)
	// run local command function if it is in the map
	localCmdFun := localCmdFunMap[inputs.Cmd]
	if localCmdFun != nil {
		if !localCmdFun(inputs) {
			return
		}
	}
	result, err := tunnel.Request(trimInput)
	if err != nil {
		tunnel.Linked = false
		fmt.Println(err.Error())
		return
	}
	fmt.Println(result)
}

func getInputs(in string) *ShellInputs {
	i := strings.Fields(in)
	c := i[0]
	lc := strings.ToLower(c)
	return &ShellInputs{
		LineTrim: in,
		RawCmd:   c,
		Cmd:      lc,
		Args:     i[1:],
	}
}

func suggest(in prompt.Document) []prompt.Suggest {
	key := in.LastKeyStroke()
	if key == prompt.Escape {
		return suggestNothing(in)
	}
	line := in.CurrentLine()
	if line == "" && key != prompt.Down {
		return suggestNothing(in)
	}
	beforeCursor := in.CurrentLineBeforeCursor()
	//lowerBefore := strings.TrimSpace(strings.ToLower(beforeCursor))
	//if i, ok := util.CmdHelpMap[lowerBefore]; ok {
	//	ch := util.CmdSuggests[i]
	//	return []prompt.Suggest{
	//		ch,
	//	}
	//}
	s := util.CmdSuggests
	return prompt.FilterHasPrefix(s, beforeCursor, true)
}

func suggestNothing(in prompt.Document) []prompt.Suggest {
	return []prompt.Suggest{}
}

func changeLivePrefix() (string, bool) {
	if tunnel.Linked {
		return tunnel.Address + "> ", true
	} else {
		return "not connected>", true
	}
}
