package form

type Users struct {
	Username    string `json: "username" `
	Online      bool   `json: "online"`
	Status_user string `json: "status_user "`
	Total_win   int    `json: "total_win"`
	Total_lose  int    `json: "total_lose"`
}

type Challenger struct {
	Username string `json: "username" `
	Win      int    `json: "win"`
	Lose     int    `json: "lose"`
}
