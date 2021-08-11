package actions

import (
	"fmt"

	"github.com/monorepo/companion_services/chimera/k8s"
	"github.com/sirupsen/logrus"
)

func BackupDatabase(env, client, slackToken, slackChannelID string) error {
	command := fmt.Sprintf("/root/scripts/db-dump --uploads3 --slack-token=%s --slack-channel-id=%s > /dev/null 2> /dev/null &", slackToken, slackChannelID)
	resp, err := k8s.Exec(env, client, k8s.ServiceCerberus, command)
	logrus.WithField("response_message", resp).Info("Exec backup db")
	return err
}
