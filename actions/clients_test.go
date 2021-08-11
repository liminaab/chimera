package actions

import (
	"reflect"
	"testing"
)

func TestFilterUnsupportedClients(t *testing.T) {
	tests := []struct {
		name    string
		clients []string
		expects []string
	}{
		{
			name:    "dev",
			clients: []string{"argocd", "cert-manager", "testA", "companions", "default", "testB", "external-secrets-controller", "ingress-nginx", "kube-node-lease", "kube-public", "kube-system", "monitoring", "testC", "testD", "testE", "service"},
			expects: []string{"testA", "testB", "testC", "testD", "testE"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FilterUnsupportedClients(tt.clients)
			if !reflect.DeepEqual(tt.expects, got) {
				t.Errorf("Expects: %s, got: %s", tt.expects, got)
			}
		})
	}
}

func TestIsUnsupportedClient(t *testing.T) {
	tests := []struct {
		name           string
		checkingClient string
		want           bool
	}{
		{
			name:           "testA",
			checkingClient: "testA",
			want:           false,
		},
		{
			name:           "argocd",
			checkingClient: "argocd",
			want:           true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsUnsupportedClient(tt.checkingClient); got != tt.want {
				t.Errorf("IsUnsupportedClient() = %v, want %v", got, tt.want)
			}
		})
	}
}
