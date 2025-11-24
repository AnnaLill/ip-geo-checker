package models

// CheckIPRequest структура для входящего запроса POST /check_ip
type CheckIPRequest struct {
	IP string `json:"ip"`
}

// CheckIPResponse структура для ответа POST /check_ip
type CheckIPResponse struct {
	City      string `json:"city"`
	FromCache bool   `json:"from_cache"`
}

// StatsResponse структура для ответа GET /get_stats
type StatsResponse struct {
	RequestsCount                   int `json:"requests_count"`
	CacheCount                      int `json:"cache_count"`
	ExternalSuccessfulRequestsCount int `json:"external_successful_requests_count"`
}
