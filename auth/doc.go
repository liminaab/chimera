package auth

/*
Full json config:
[
	{
		slack_user_id: "1",
		permissions: {
			BackupDB: ["dev", "stage", "prod"],
			RestartService: ["dev", "stage"]
		},
	},
	{
		slack_user_id: "2",
		permissions: {
			BackupDB: ["*"],
			RestartService: ["*"]
		},
	},
]
*/
