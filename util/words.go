package util

const (
	Version    = "1.0.0"
	ShellTitle = "Redis client"
	CmdUsage   = "redis-cli"
	CmdShort   = "redis client to connect and manage server"
	CmdLong    = ""
	CmdExample = `
  cat /etc/passwd | redis-client -x set mypasswd
  redis-client get mypasswd
  redis-client -r 100 lpush mylist x
  redis-client -r 100 -i 1 info | grep used_memory_human:
  redis-client --eval myscript.lua key1 key2 , arg1 arg2 arg3
  redis-client --scan --pattern '*:12345*'

  (Note: when using --eval the comma separates KEYS[] from ARGV[] items)

When no command is given, redis-cli starts in interactive mode.
Type "help" in interactive mode for information on available commands
and settings.
`
	ShellHelp = CmdUsage + " " + Version + `
To get help about Redis commands type:
      "help @<group>" to get a list of commands in <group>
      "help <command>" for help on <command>
      "help <tab>" to get a list of possible help topics
      "quit" to exit

To set redis-cli preferences:
      ":set hints" enable online hints
      ":set nohints" disable online hints
Set your preferences in ~/.redisclientrc
`
)
