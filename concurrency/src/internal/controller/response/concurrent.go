package response

type ConcurrentResp struct {
	ID       string `json:"id"`
	Status   string `json:"status"`
	Contents string `json:"contents,omitempty"`
}
