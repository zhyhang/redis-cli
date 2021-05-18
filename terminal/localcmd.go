package terminal

import "strings"

var exit = map[string]bool{
	"quit": true,
	"exit": true,
}

func exitSignal(cmd string) bool {
	lowerCmd := strings.ToLower(cmd)
	if exit[lowerCmd] {
		return true
	}
	return false
}
