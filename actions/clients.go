package actions

import (
	"sort"

	"github.com/monorepo/companion_services/chimera/k8s"
)

var unsupportedClients = []string{
	"argocd",
	"cert-manager",
	"companions",
	"default",
	"external-secrets-controller",
	"ingress-nginx",
	"kube-node-lease",
	"kube-public",
	"kube-system",
	"monitoring",
	"service",
}

func GetClients(environment string) ([]string, bool, error) {
	if valid := k8s.IsValidNameSpace(environment); !valid {
		return []string{}, false, nil
	}

	ns, err := k8s.GetNamespaces(environment)
	if err != nil {
		return []string{}, false, err
	}
	ns = FilterUnsupportedClients(ns)
	sort.Strings(ns)
	return ns, true, nil
}

func FilterUnsupportedClients(clients []string) (newClient []string) {
	if len(clients) == 0 {
		return
	}
	for i := len(clients) - 1; i >= 0; i-- {
		client := clients[i]
		if !IsUnsupportedClient(client) {
			continue
		}
		clients = append(clients[:i], clients[i+1:]...)
	}
	return clients
}

func IsUnsupportedClient(checkingClient string) bool {
	for _, unsupportedClient := range unsupportedClients {
		if unsupportedClient == checkingClient {
			return true
		}
	}
	return false
}
