package auth

type PermissionName string

const (
	PermissionNameBackupDB       = "BackupDB"
	PermissionNameRestartService = "RestartService"
)

type EnvironmentName string

type Permissions map[PermissionName]Permission

type Permission []string

func (p Permission) ValidToDo(pattern string) bool {
	for _, item := range p {
		if item == pattern || item == "*" {
			return true
		}
	}
	return false
}
