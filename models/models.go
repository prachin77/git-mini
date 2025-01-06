package models

type UserConfig struct {
	Username       string                `json:"username"`
	Port           string                `json:"port"`
	SendWorkSpaces []SendWorkSpaceFolder `json:"send_workspaces"`

	// later when dynamic IPs are stored , bcz right now we're manually entering IPs
	// Ip_Address     string                `json:"ip_address"`
}

// workspace hosted by a user to send to connected clients who are retrieving data from the folder 
type SendWorkSpaceFolder struct {
	Workspace_Name        string `json:"workspace_name"`
	Workspace_Password    string `json:"workspace_password"`
	Workspace_Path        string `json:"workspace_path"`
	Workspace_Hosted_Date string `json:"workspace_hosted_date"`
}

type Files struct {
	FileName     string `json:"file_name"`
	FileSize     string `json:"file_size"`
	FileLocation string `json:"file_location"`
}
