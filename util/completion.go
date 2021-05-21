package util

import "github.com/c-bata/go-prompt"

var CmdSuggests = createCmdSuggests()

func createCmdSuggests() []prompt.Suggest {
	helps := GetCommandHelps()
	length := len(helps)
	ss := make([]prompt.Suggest, length)
	for i := 0; i < length; i++ {
		ss[i] = prompt.Suggest{
			Text:        helps[i].Name,
			Description: helps[i].Params,
		}
	}
	return ss
}
