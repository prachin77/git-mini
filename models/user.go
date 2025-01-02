package models

type UserConfig struct {
	Username           string `json:"username"`
	Workspace_Password string `json:"workspace_password"`
	Port               string `json:"port"`
	Ip_Address         string `json:"ip_address"`
}
