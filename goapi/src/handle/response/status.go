package response

const (
	COMPLETED = "COMPLETED"
	FAILED    = "FAILED"
	RUNNING   = "RUNNING"
	PENDING   = "PENDING"
)

type StatusRes struct {
	Id     string `json:"id"`
	Status string `json:"status"`
}
