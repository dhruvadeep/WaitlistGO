package models

type Ping struct {
	Response int    `json:"response"`
	Message  string `json:"message"`
	IP       string `json:"ip"`
}