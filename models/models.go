package models

var StatusInfo map[string]string

func init() {
	StatusInfo = make(map[string]string)
}

type PostRequest struct {
	AddedWebsites []string `json:"websites"`
}
