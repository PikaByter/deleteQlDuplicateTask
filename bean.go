package main

type loginResp struct {
	Code	int
	Token	string
}

type Task struct {
	Command string `json:"command"`
	Created int	`json:"created"`
	IsDisabled int `json:"isDisabled"`
	IsSystem int `json:"isSystem"`
	LogPath string `json:"log_path"`
	Name string `json:"name"`
	Saved bool `json:"saved"`
	Schedule string `json:"schedule"`
	Status int `json:"status"`
	Timestamp string `json:"timestamp"`
	Id string `json:"_id"`
}

type getTaskResp struct{
	Code int `json:"code"`
	Data []Task `json:"data"`
}

type trashResp struct {
	Code int `json:"code"`
}