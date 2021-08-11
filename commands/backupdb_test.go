package commands

import "testing"

func Test_parseBackupDBArgs(t *testing.T) {
	tests := []struct {
		name       string
		raw        string
		wantClient string
		wantEnv    string
	}{
		{
			name:       "qa.dev",
			raw:        "qa.dev",
			wantClient: "qa",
			wantEnv:    "dev",
		},
		{
			name:       "pr1271.test",
			raw:        "pr1271.test",
			wantClient: "pr1271",
			wantEnv:    "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotClient, gotEnv := parseBackupDBArgs(tt.raw)
			if gotClient != tt.wantClient {
				t.Errorf("parseBackupDBArgs() gotClient = %v, want %v", gotClient, tt.wantClient)
			}
			if gotEnv != tt.wantEnv {
				t.Errorf("parseBackupDBArgs() gotEnv = %v, want %v", gotEnv, tt.wantEnv)
			}
		})
	}
}
