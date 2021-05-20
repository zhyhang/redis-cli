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
		suggestNothing,
		prompt.OptionPrefix(prefix),
		prompt.OptionLivePrefix(changeLivePrefix),
		prompt.OptionTitle(util.ShellTitle),
		prompt.OptionMaxSuggestion(1),
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

func suggestNothing(in prompt.Document) []prompt.Suggest {
	return []prompt.Suggest{}
}

func suggest(in prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "users", Description: "Store the username and age"},
		{Text: "articles", Description: "Store the article text posted by user"},
		{Text: "comments", Description: "Store the text commented to articles"},
		{Text: "groups", Description: "Combine users with specific rules"},
	}
	return prompt.FilterHasPrefix(s, in.GetWordBeforeCursor(), true)
}

func changeLivePrefix() (string, bool) {
	if tunnel.Linked {
		return tunnel.Address + "> ", true
	} else {
		return "not connected>", true
	}
}
