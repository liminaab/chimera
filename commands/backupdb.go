package commands

import (
	"fmt"
	"strings"

	"github.com/monorepo/companion_services/chimera/actions"
	"github.com/monorepo/companion_services/chimera/auth"
	"github.com/monorepo/companion_services/chimera/slackutils"
)

func BackupDB(slackRequest *slackutils.MessageEventRequest, args string) error {
	client, env := parseBackupDBArgs(args)
	if !actions.IsValidEnvironment(env) {
		return slackutils.SendSimple(slackutils.GetOAuthToken(),
			slackRequest.Event.Channel, fmt.Sprintf("We don't support environment '%s'.", env))
	}
	if actions.IsUnsupportedClient(client) {
		return slackutils.SendSimple(slackutils.GetOAuthToken(),
			slackRequest.Event.Channel, fmt.Sprintf("We don't support client '%s'.", client))
	}
	if !auth.IsValidToBackupDB(slackRequest.Event.User, env) {
		return slackutils.SendSimple(slackutils.GetOAuthToken(),
			slackRequest.Event.Channel, fmt.Sprintf("You don't have permission to backup db on '%s'.", env))
	}

	go slackutils.SendSimple(slackutils.GetOAuthToken(), slackRequest.Event.Channel, "The backup is processing. I will send you download link asap.")
	return actions.BackupDatabase(env, client, slackutils.GetOAuthToken(), slackRequest.Event.Channel)
}

func parseBackupDBArgs(raw string) (client, env string) {
	raw = strings.TrimPrefix(raw, "<http://") // trick because sometime slack autofills
	args := strings.SplitN(raw, ".", 2)
	if len(args) < 2 {
		return raw, ""
	}
	return args[0], args[1]
}
