package api

type EngageResponse struct {
	Page      int    `json:"page"`
	PageSize  int    `json:"page"`
	SessionID string `json:"session_id"`
	Status    string `json:"status"`
	Total     int    `json:"total"`

	Results []map[string]interface{} `json:"results"`
}
