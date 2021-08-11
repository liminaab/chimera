package k8s

import (
	"encoding/json"
	"os"

	"github.com/pkg/errors"
)

/*
[
	{
		"name": "test",
		"endpoint": "https://ABVC.yl4.eu-north-1.eks.amazonaws.com",
		"certificate_authority": "LS0tLS1CRUdJTiBDRVJUSUZJQ0..."
	}
]
*/
type Config struct {
	Name                 string `json:"name"`
	Endpoint             string `json:"endpoint"`
	CertificateAuthority string `json:"certificate_authority"`
}

type Configs []Config

func ParseClusterConfigs(raw []byte) (configs Configs, err error) {
	return configs, json.Unmarshal(raw, &configs)
}

func ParseAuthsFromEnv() (configs Configs, err error) {
	return ParseClusterConfigs([]byte(os.Getenv("CLUSTERS")))
}

func (cs Configs) GetByName(name string) (Config, error) {
	for _, c := range cs {
		if c.Name == name {
			return c, nil
		}
	}
	return Config{}, errors.Errorf("Cluster '%s' doesn't supported")
}

func (cs Configs) GetNames() (names []string) {
	for _, c := range cs {
		names = append(names, c.Name)
	}
	return names
}
