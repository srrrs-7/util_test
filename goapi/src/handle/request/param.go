package request

type Params struct {
	UserId string `json:"user_id"`
	From   string `json:"from"`
	To     string `json:"to"`
}
