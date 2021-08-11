package commands

import (
	"fmt"
	"strings"

	"github.com/monorepo/companion_services/chimera/actions"
	"github.com/monorepo/companion_services/chimera/slackutils"
)

func ListEnvironments(slackRequest *slackutils.MessageEventRequest) error {
	environments := actions.GetAllowedEnvironments()
	if len(environments) == 0 {
		return slackutils.SendSimple(slackutils.GetOAuthToken(),
			slackRequest.Event.Channel, fmt.Sprintf("Currently, we don't support any client"))
	}
	envStr := strings.Join(environments, "\n")
	return slackutils.SendSimple(slackutils.GetOAuthToken(),
		slackRequest.Event.Channel, fmt.Sprintf("System supports environments: ```%s```", envStr))
}

func ListClients(slackRequest *slackutils.MessageEventRequest, env string) error {
	clients, existed, err := actions.GetClients(env)
	if err != nil {
		return err
	}
	if !existed {
		return slackutils.SendSimple(slackutils.GetOAuthToken(),
			slackRequest.Event.Channel, fmt.Sprintf("Can't query clients. The environment '%s' doest existed", env))
	}

	if len(clients) == 0 {
		return slackutils.SendSimple(slackutils.GetOAuthToken(),
			slackRequest.Event.Channel, fmt.Sprintf("Currently, '%s' doesn't have any client", env))
	}

	for index := range clients {
		clients[index] = fmt.Sprintf("%s.%s", clients[index], env)
	}
	clientStr := strings.Join(clients, "\n")
	return slackutils.SendSimple(slackutils.GetOAuthToken(),
		slackRequest.Event.Channel, fmt.Sprintf("'%s' supports: ```%s```", env, clientStr))
}
