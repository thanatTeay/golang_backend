package form

type Challenge struct {
	Username          string `json: "username"`
	Challenger        string `json: "challenger"`
	Action_user       string `json: "action_user"`
	Action_challenger string `json: "action_challenger"`
}
