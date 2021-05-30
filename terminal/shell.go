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
	before := strings.TrimSpace(in.CurrentLineBeforeCursor())
	beforeWord := in.GetWordBeforeCursor()
	css := util.CmdSuggests
	bss := prompt.FilterHasPrefix(css, before, true)
	if before == beforeWord {
		return bss
	}
	beforeTrim := strings.TrimSuffix(before, beforeWord)
	bss = prompt.FilterHasPrefix(css, beforeTrim, true)
	rss := make([]prompt.Suggest, 0, len(bss))
	for _, s := range bss {
		newText := strings.TrimSpace(strings.TrimPrefix(s.Text, strings.ToUpper(beforeTrim)))
		if newText != "" {
			rss = append(rss, prompt.Suggest{newText, s.Description})
		}
	}
	return prompt.FilterHasPrefix(rss, beforeWord, true)
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
