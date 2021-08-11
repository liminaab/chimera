package commands

import (
	"github.com/monorepo/companion_services/chimera/slackutils"
	"github.com/slack-go/slack"
)

const helpMsg = `Can I help you? I support commands:
 - Backup database: ` + "`" + `backupdb ${client}.${env_name}` + "`" + `. e.g. ` + "`" + `backupdb qa.dev` + "`" + `
 - Show environments: ` + "`" + `envs` + "`" + `
 - Show clients: ` + "`" + `envs ${env_name}` + "`" + `. e.g. ` + "`" + `envs dev` + "`" + `
`

func RootHelp(slackRequest *slackutils.MessageEventRequest) error {
	return slackutils.SendComplex(slackutils.GetOAuthToken(), slackRequest.Event.Channel,
		slack.MsgOptionCompose(
			slack.MsgOptionText(helpMsg, true),
		),
	)
}

const invalidMsg = "Your command is invalid. Please use `help` to check again!"

func InvalidCommand(slackRequest *slackutils.MessageEventRequest) error {
	return slackutils.SendSimple(
		slackutils.GetOAuthToken(),
		slackRequest.Event.Channel,
		invalidMsg,
	)
}
