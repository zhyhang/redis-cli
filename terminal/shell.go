package terminal

import (
	"fmt"
	prompt "github.com/c-bata/go-prompt"
	"github.com/zhyhang/redis-client/redis"
	"github.com/zhyhang/redis-client/util"
	"strings"
)

var tunnel *redis.Tunnel
var config = newShellConfig()

func Run(flags *CmdFlags) {
	config.cmdLine = flags
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
	if isIgnoreInput() {
		return
	}
	trimInput := strings.TrimSpace(input)
	if trimInput == "" {
		return
	}
	inputs := getInputs(trimInput)
	// update config by input info
	updateConf(inputs, config)
	// run local command function if it is in the map
	localCmdFun := localCmdFunMap[inputs.Cmd]
	if localCmdFun != nil {
		if !localCmdFun(inputs) {
			return
		}
	}
	result, err := tunnel.Request(trimInput)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(result)
	if config.modeMonitor {
		go exeKeepaliveCmd()
	}
}

func exeKeepaliveCmd() {
	for {
		reading, err := tunnel.KeepReading()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Println(reading)
	}
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
	if isIgnoreInput() {
		return suggestNothing(in)
	}
	key := in.LastKeyStroke()
	if key == prompt.Escape {
		return suggestNothing(in)
	}
	if in.Text == "" && key != prompt.Down {
		return suggestNothing(in)
	}
	localSs := suggestLocalCmd(in)
	if localSs != nil {
		return localSs
	}
	before := in.TextBeforeCursor()
	beforeWord := in.GetWordBeforeCursor()
	css := util.CmdSuggests
	// input only contains one word
	if before == beforeWord {
		return prompt.FilterHasPrefix(css, before, true)
	}
	beforeTrim := strings.TrimSuffix(before, beforeWord)
	// already input whole command
	if cmdIdx, ok := util.CmdHelpMap[strings.TrimSpace(strings.ToLower(beforeTrim))]; ok {
		return suggestCmdPara(0, util.CmdHelps[cmdIdx])
	}
	// already input whole command and some parameters
	if ss, found := suggestCanFilterCmd(in); found {
		return ss
	}
	// command contains multi words
	pss := prompt.FilterHasPrefix(css, beforeTrim, true)
	if len(pss) > 0 {
		nss := make([]prompt.Suggest, 0, len(pss))
		for _, ps := range pss {
			textTrim := strings.TrimPrefix(ps.Text, strings.ToUpper(beforeTrim))
			if textTrim == "" {
				nss = append(nss, prompt.Suggest{" ", ps.Description})
			} else {
				nss = append(nss, prompt.Suggest{strings.TrimSpace(textTrim), ps.Description})
			}
		}
		return prompt.FilterHasPrefix(nss, beforeWord, true)
	} else {
		return suggestNothing(in)
	}
}

func suggestCanFilterCmd(in prompt.Document) (ss []prompt.Suggest, found bool) {
	after := in.TextAfterCursor()
	beforeIncludeCur := in.TextBeforeCursor()
	if after != "" {
		beforeIncludeCur = in.TextBeforeCursor() + (after[0:1])
	}
	fields := strings.Fields(beforeIncludeCur)
	for i := 1; i < len(fields); i++ {
		cmd := strings.ToLower(strings.Join(fields[0:i], " "))
		if cmdIdx, ok := util.CmdHelpMap[cmd]; ok {
			found = true
			ss = suggestCmdPara(len(fields)-i-1, util.CmdHelps[cmdIdx])
			return
		}
	}
	ss = nil
	found = false
	return
}

func suggestLocalCmd(in prompt.Document) (ss []prompt.Suggest) {
	fields := strings.Fields(in.TextBeforeCursor())
	localCmd := strings.ToLower(fields[0])
	if localCmd == "help" {
		if len(fields) == 1 && strings.HasSuffix(in.TextBeforeCursor(), " ") {
			return util.LocalHelpCmdSuggests
		} else if len(fields) == 2 && !strings.HasSuffix(in.TextBeforeCursor(), " ") {
			return prompt.FilterHasPrefix(util.LocalHelpCmdSuggests, fields[1], true)
		} else {
			return suggestNothing(in)
		}
	}
	return nil
}

func suggestCmdPara(curParaIdx int, cmdHelp util.CommandHelp) []prompt.Suggest {
	return []prompt.Suggest{{" ", cmdHelp.Params}}
	//ss := []prompt.Suggest{{" ", cmdHelp.Params}}
	//if cmdHelp.Params == "-" {
	//	return ss
	//}
	//paraFields := strings.Fields(cmdHelp.Params)
	//hint := ""
	//for i := 0; i < curParaIdx; i++ {
	//	fl := len(paraFields[i])
	//	hint += strings.Repeat("~", int(math.Min(3, float64(fl)))) + " "
	//}
	//fl := len(paraFields[curParaIdx])
	//hint += strings.Repeat(".", int(math.Min(3, float64(fl))))
	//return append(ss, prompt.Suggest{
	//	" ",
	//	hint,
	//})
}

func suggestNothing(in prompt.Document) []prompt.Suggest {
	return []prompt.Suggest{}
}

func changeLivePrefix() (string, bool) {
	if isIgnoreInput() {
		return "", true
	}
	if tunnel.Linked {
		return tunnel.Address + "> ", true
	} else {
		return "not connected>", true
	}
}

func isIgnoreInput() bool {
	return config.modeMonitor
}
