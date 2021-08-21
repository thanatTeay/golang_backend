package models

type Users struct {
	Username string `json: "username" `
	Online   bool   `json: "online"`
}
