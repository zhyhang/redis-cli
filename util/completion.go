package util

import "github.com/c-bata/go-prompt"

var CmdSuggests = createCmdSuggests()
var LocalHelpCmdSuggests = createHelpLocalCmdSuggests()

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

func createHelpLocalCmdSuggests() []prompt.Suggest {
	groups := GetCommandGroups()
	ss := make([]prompt.Suggest, 0, len(groups)+len(CmdSuggests))
	for _, g := range groups {
		ss = append(ss, prompt.Suggest{
			"@" + g,
			"",
		})
	}
	for _, s := range CmdSuggests {
		ss = append(ss, s)
	}
	return ss
}
