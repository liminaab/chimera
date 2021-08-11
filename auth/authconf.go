package auth

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Auth struct {
	SlackUserID string      `json:"slack_user_id"`
	Permissions Permissions `json:"permissions"`
}

type Auths []Auth

func (auths Auths) GetPermissionsBySlackID(checkingSlackID string) (Permissions, error) {
	for _, auth := range auths {
		slackUserIDs := strings.Split(auth.SlackUserID, ",")
		for _, slackUserID := range slackUserIDs {
			if slackUserID == checkingSlackID {
				return auth.Permissions, nil
			}
		}
	}
	return nil, errors.Errorf("SlackUserID %s doesn't definied", checkingSlackID)
}

func ParseAuths(raw []byte) (auths Auths, err error) {
	return auths, json.Unmarshal(raw, &auths)
}

func ParseAuthsFromEnv() (auths Auths, err error) {
	return ParseAuths([]byte(os.Getenv("AUTHS")))
}

func IsValidToBackupDB(slackID, environment string) bool {
	auths, err := ParseAuthsFromEnv()
	if err != nil {
		logrus.Warn(err)
		return false
	}

	permissions, err := auths.GetPermissionsBySlackID(slackID)
	if err != nil {
		logrus.Warn(err)
		return false
	}

	permission, _ := permissions[PermissionNameBackupDB]
	return permission.ValidToDo(environment)
}
