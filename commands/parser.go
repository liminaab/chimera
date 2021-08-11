package commands

import (
	"strings"

	"github.com/monorepo/companion_services/chimera/slackutils"
)

type CommandIntent int

const (
	IntentRootHelp         CommandIntent = iota
	IntentBackupDB         CommandIntent = iota
	IntentListEnvironments CommandIntent = iota
	IntentListClients      CommandIntent = iota
	IntentInvalidCommand   CommandIntent = iota
)

type Command string

const (
	CommandHi               Command = "hi"
	CommandHello            Command = "hello"
	CommandHej              Command = "hej"
	CommandHelp             Command = "help"
	CommandBackupDB         Command = "backupdb"
	CommandListEnvironments Command = "envs"
	CommandListEnvironment  Command = "env"
)

func ParseIntent(request *slackutils.MessageEventRequest) (CommandIntent, string) {
	args := strings.Split(strings.TrimSpace(request.Event.Text), " ")
	if len(args) == 0 {
		return IntentInvalidCommand, ""
	}

	command := Command(strings.ToLower(args[0]))
	argsStr := strings.Join(args[1:], "")

	switch command {
	case CommandHi, CommandHello, CommandHej, CommandHelp:
		return IntentRootHelp, argsStr
	case CommandBackupDB:
		return IntentBackupDB, argsStr
	case CommandListEnvironments, CommandListEnvironment:
		if len(argsStr) > 0 { // Has parameter to query clients of an environment
			return IntentListClients, argsStr
		} else {
			return IntentListEnvironments, argsStr
		}
	default:
		return IntentInvalidCommand, argsStr
	}
}
