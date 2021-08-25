package models

type Challenge struct {
	Id string `json: "id"`
	//Users1            Users  `json: "users1"`
	//Users2            Users  `json: "users2"`
	User              string `json: "user"`
	Challenger        string `json: "challenger"`
	User_win          int    `json: "user_win"`
	Challenger_win    int    `json: "challenger_win"`
	User_lose         int    `json: "user_lose"`
	Challenger_lose   int    `json: "challenger_lose"`
	Action_user1      string `json: "action_user"`
	Action_challenger string `json: "action_challenger"`
	Type              int    `json: "type"`
	Winner            int    `json: "winner"`
}

type History struct {
	Id                string `json: "id"`
	Date              string `json: "date"`
	User              string `json: "user"`
	Challenger        string `json: "challenger"`
	User_win          int    `json: "user_win"`
	Challenger_win    int    `json: "challenger_win"`
	User_lose         int    `json: "user_lose"`
	Challenger_lose   int    `json: "challenger_lose"`
	Action_user       string `json: "action_user"`
	Action_challenger string `json: "action_challenger"`
}
