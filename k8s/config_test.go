package k8s

import (
	"os"
	"reflect"
	"testing"
)

func TestParseAuthsFromEnv(t *testing.T) {
	os.Setenv("CLUSTERS", `[
{
	"name": "test",
	"endpoint": "https://ABC.com",
	"certificate_authority": "ABC"
},
{
	"name": "dev",
	"endpoint": "https://XYZ.com",
	"certificate_authority": "XYZ"
}
]`)
	tests := []struct {
		name        string
		wantConfigs Configs
		wantErr     bool
	}{
		{
			name: "Valid",
			wantConfigs: Configs{
				{
					Name:                 "test",
					Endpoint:             "https://ABC.com",
					CertificateAuthority: "ABC",
				},
				{
					Name:                 "dev",
					Endpoint:             "https://XYZ.com",
					CertificateAuthority: "XYZ",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotConfigs, err := ParseAuthsFromEnv()
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseClusterConfigs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotConfigs, tt.wantConfigs) {
				t.Errorf("ParseClusterConfigs() = %v, want %v", gotConfigs, tt.wantConfigs)
			}
		})
	}
}
