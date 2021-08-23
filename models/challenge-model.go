package models

import "time"

type Challenge struct {
	Id           string `json: "id"`
	Users1       Users
	Users2       Users
	Action_user1 string `json: "action_user"`
	Action_user2 string `json: "action_challenger"`
}

type History struct {
	Date              time.Time `json: "date"`
	User1             string    `json "user1"`
	User2             string    `json "user2"`
	User1_win         string    `json "user1_win"`
	User2_win         string    `json "user2_win"`
	User1_lose        string    `json "user1_lose"`
	User2_lose        string    `json "user2_lose"`
	Action_user       string    `json: "action_user"`
	Action_challenger string    `json: "action_challenger"`
}
