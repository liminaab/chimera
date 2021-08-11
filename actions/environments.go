package actions

import (
	"github.com/monorepo/companion_services/chimera/k8s"
	"github.com/sirupsen/logrus"
)

func GetAllowedEnvironments() []string {
	configs, err := k8s.ParseAuthsFromEnv()
	if err != nil {
		logrus.WithError(err).Error("Get allowed environments")
	}
	allowedEnvironments := make([]string, 0, len(configs))
	for _, config := range configs {
		allowedEnvironments = append(allowedEnvironments, config.Name)
	}
	return allowedEnvironments
}

func IsValidEnvironment(env string) bool {
	return k8s.IsValidNameSpace(env)
}
