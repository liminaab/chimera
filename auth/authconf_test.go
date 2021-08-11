package auth

import (
	"os"
	"reflect"
	"testing"
)

func TestParseAuth(t *testing.T) {
	authRaw := `[
   {
      "slack_user_id":"1",
      "permissions":{
         "BackupDB":[
            "dev",
            "stage",
            "prod"
         ],
         "RestartService":[
            "dev",
            "stage"
         ]
      }
   },
   {
      "slack_user_id":"2",
      "permissions":{
         "BackupDB":[
            "*"
         ],
         "RestartService":[
            "*"
         ]
      }
   }
]`
	gotAuths, err := ParseAuths([]byte(authRaw))
	if err != nil {
		t.Errorf("got error: %v", err)
		return
	}

	expectedLength := 2
	if len(gotAuths) != expectedLength {
		t.Errorf("expected length: %v, got length: %v", expectedLength, len(gotAuths))
		return
	}

	expects := Auths{
		{
			SlackUserID: "1",
			Permissions: Permissions{
				"BackupDB":       []string{"dev", "stage", "prod"},
				"RestartService": []string{"dev", "stage"},
			},
		},
		{
			SlackUserID: "2",
			Permissions: Permissions{
				"BackupDB":       []string{"*"},
				"RestartService": []string{"*"},
			},
		},
	}

	for index := range gotAuths {
		got := gotAuths[index]
		expect := expects[index]
		if !reflect.DeepEqual(got, expect) {
			t.Errorf("expected: %+v, got: %+v", expect, got)
			return
		}
	}
}

func TestIsValidToBackupDB(t *testing.T) {
	os.Setenv("AUTHS", `[
   {
      "slack_user_id":"1",
      "permissions":{
         "BackupDB":[
            "dev",
            "stage",
            "prod"
         ],
         "RestartService":[
            "dev",
            "stage"
         ]
      }
   },
   {
      "slack_user_id":"2",
      "permissions":{
         "BackupDB":[
            "*"
         ],
         "RestartService":[
            "*"
         ]
      }
   }
]`)
	type args struct {
		slackID     string
		environment string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "valid",
			args: args{
				slackID:     "1",
				environment: "stage",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidToBackupDB(tt.args.slackID, tt.args.environment); got != tt.want {
				t.Errorf("IsValidToBackupDB() = %v, want %v", got, tt.want)
			}
		})
	}
}
