package slackutils

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
)

func GetOAuthToken() string {
	return os.Getenv("SLACK_OAUTH_TOKEN")
}

func SendSimple(token, channelID, message string) error {
	msgOption := slack.MsgOptionCompose(
		slack.MsgOptionText(message, true),
	)
	respChannel, respTimestamp, err := slack.New(token).PostMessage(channelID, msgOption)
	logrus.WithFields(logrus.Fields{
		"resp_channel":   respChannel,
		"resp_timestamp": respTimestamp,
		"error":          err,
		"token":          token,
		"channel_id":     channelID,
		"message":        message,
	}).Info("SendSlackMessage simple")
	return err
}

func SendComplex(token, channelID string, options slack.MsgOption) error {
	respChannel, respTimestamp, err := slack.New(token).PostMessage(channelID, options)
	logrus.WithFields(logrus.Fields{
		"resp_channel":   respChannel,
		"resp_timestamp": respTimestamp,
		"error":          err,
		"token":          token,
		"channel_id":     channelID,
	}).Info("SendSlackMessage complex")
	return err
}
