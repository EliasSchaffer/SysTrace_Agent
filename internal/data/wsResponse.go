package data

type WSResponse struct {
	Type    string                 `json:"type"`
	Status  string                 `json:"status,omitempty"`
	Message string                 `json:"message,omitempty"`
	Payload map[string]interface{} `json:"payload,omitempty"`
}
