package models

type UserConfig struct {
	Username           string                    `json:"username"`
	Port               string                    `json:"port"`
	SendWorkSpaces     []SendWorkSpaceFolder     `json:"send_workspaces"`
	RecievedWorkspaces []RecievedWorkSpaceFolder `json:"recieved workspaces"`
	// later when dynamic IPs are stored , bcz right now we're manually entering IPs
	// Ip_Address     string                `json:"ip_address"`
}

type SendWorkSpaceFolder struct {
	Workspace_Name        string `json:"workspace_name"`
	Workspace_Password    string `json:"workspace_password"`
	Workspace_Path        string `json:"workspace_path"`
	Workspace_Hosted_Date string `json:"workspace_hosted_date"`
}

type RecievedWorkSpaceFolder struct {
	Workspace_Name string `json:"workspace_name"`
	Workspace_Path string `json:"workspace_path"`
	Workspace_IP   string `json:"workspace_ip"`
	Recieved_Date  string `json:"recieved_date"`
	// analyse its usage first & then use
	// LastHash string `json:"last_hast"`
}

type Files struct {
	FileName     string `json:"file_name"`
	FileSize     string `json:"file_size"`
	FileLocation string `json:"file_location"`
}

type WorkspaceConfig struct {
	WorkspaceName     string       `json:"workspace_name"`
	WorkspacePath    string       `json:"workspace_path"`
	Hosted_By         string       `json:"hosted_by"`
	Port              string       `json:"port"`
	WorkspacePassword string       `json:"workspace_password"`
	No_Of_Files       int          `json:"no_of_files"`
	Workspace_Size    string       `json:"workspace_size"`
	InitDate          string       `json:"init_date"`
	File              []Files      `json:"files"`
	Connections       []Connection `json:"connections"`
}

// people who've cloned/connected (to)the workspace
type Connection struct {
	Username  string `json:"username"`
	Port      string `json:"port"`
	CurrentIp string `json:"current_ip"`
}