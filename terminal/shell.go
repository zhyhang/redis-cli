package terminal

import (
	"fmt"
	prompt "github.com/c-bata/go-prompt"
	"github.com/zhyhang/redis-client/platform"
	"github.com/zhyhang/redis-client/redis"
	"os"
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
		prompt.OptionTitle("Redis client"),
		prompt.OptionMaxSuggestion(1),
		prompt.OptionPrefixTextColor(prompt.Green),
	)
	p.Run()
}

func Exit() {
	err := tunnel.Destroy()
	if err != nil {
		os.Exit(1)
	}
	platform.HandleExit()
	os.Exit(0)
}

func exec(input string) {
	trimInput := strings.TrimSpace(input)
	if trimInput == "" {
		return
	}
	if exitSignal(getCmd(trimInput)) {
		Exit()
	}
	result, err := tunnel.Request(trimInput)
	if err != nil {
		tunnel.Linked = false
		fmt.Println(err.Error())
		return
	}
	fmt.Println(result)
}

func getCmd(in string) string {
	return strings.Split(in, " ")[0]
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
