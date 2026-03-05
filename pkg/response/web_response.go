package response

type WebResponse struct {
	Code   int    `json:"code"`
	Status string `json:"status"`
	Data   any    `json:"data,omitempty"`
	Errors any    `json:"errors,omitempty"`
}
