package common

type Metadata struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type StandardResponse struct {
	Meta Metadata    `json:"meta"`
	Data interface{} `json:"data,omitempty"`
}
