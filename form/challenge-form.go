package form

type Challenge struct {
	Username          string `json: "username"`
	Challenger        string `json: "challenger"`
	Action_user       string `json: "action_user"`
	Action_challenger string `json: "action_challenger"`
	User1             Users
	User2             Users
}

type Ranking struct {
	User            string `json: "user"`
	Challenger      string `json: "challenger"`
	User_win        int    `json: "user_win"`
	Challenger_win  int    `json: "challenger_win"`
	User_lose       int    `json: "user_lose"`
	Challenger_lose int    `json: "challenger_lose"`
	Action_user1    string `json: "action_user"`
	Winner          int    `json: "winner"`
}
