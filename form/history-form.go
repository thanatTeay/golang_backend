package form

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
