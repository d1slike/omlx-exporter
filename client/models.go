package client

type StatsResponse struct {
	TotalPromptTokens     float64              `json:"total_prompt_tokens"`
	TotalCompletionTokens float64              `json:"total_completion_tokens"`
	TotalCachedTokens     float64              `json:"total_cached_tokens"`
	TotalRequests         float64              `json:"total_requests"`
	CacheEfficiency       float64              `json:"cache_efficiency"`
	AvgPrefillTPS         float64              `json:"avg_prefill_tps"`
	AvgGenerationTPS      float64              `json:"avg_generation_tps"`
	UptimeSeconds         float64              `json:"uptime_seconds"`
	ActiveModels          ActiveModelsResponse `json:"active_models"`
	RuntimeCache          RuntimeCacheResponse `json:"runtime_cache"`
}

type ActiveModelsResponse struct {
	Models             []ActiveModel          `json:"models"`
	ModelMemoryUsed    float64                `json:"model_memory_used"`
	ModelMemoryMax     float64                `json:"model_memory_max"`
	MemoryPressure     MemoryPressureResponse `json:"memory_pressure"`
	TotalActiveRequests float64               `json:"total_active_requests"`
	TotalWaitingRequests float64              `json:"total_waiting_requests"`
}

type MemoryPressureResponse struct {
	Enabled       bool   `json:"enabled"`
	CurrentBytes  float64 `json:"current_bytes"`
	SoftBytes     float64 `json:"soft_bytes"`
	HardBytes     float64 `json:"hard_bytes"`
	PressureLevel string  `json:"pressure_level"`
}

type ActiveModel struct {
	ModelID        string           `json:"id"`
	ActiveRequests int              `json:"active_requests"`
	WaitingRequests int             `json:"waiting_requests"`
	Prefilling     []PrefillRequest `json:"prefilling"`
	Generating     []RequestDetail  `json:"generating"`
	Waiting        []WaitingRequest `json:"waiting"`
	IdleSeconds    float64          `json:"idle_seconds"`
}

type RequestDetail struct {
	RequestID              string  `json:"request_id"`
	ElapsedSeconds         float64 `json:"elapsed_seconds"`
	GeneratedTokens        int     `json:"generated_tokens"`
	TokensPerSecond        float64 `json:"tokens_per_second"`
	LastActivityAgeSeconds float64 `json:"last_activity_age_seconds"`
	PromptTokens           int     `json:"prompt_tokens"`
	MaxTokens              int     `json:"max_tokens"`
}

type PrefillRequest struct {
	RequestID string  `json:"request_id"`
	Processed int     `json:"processed"`
	Total     int     `json:"total"`
	Speed     float64 `json:"speed"`
	ETA       float64 `json:"eta"`
	Elapsed   float64 `json:"elapsed"`
	Phase     string  `json:"phase"`
}

type WaitingRequest struct {
	RequestID      string  `json:"request_id"`
	QueuePosition  int     `json:"queue_position"`
	ElapsedSeconds float64 `json:"elapsed_seconds"`
	PromptTokens   int     `json:"prompt_tokens"`
}

type RuntimeCacheResponse struct {
	TotalNumFiles     int     `json:"total_num_files"`
	TotalSizeBytes    float64 `json:"total_size_bytes"`
	DiskMaxBytes      float64 `json:"disk_max_bytes"`
	HotCacheSizeBytes float64 `json:"hot_cache_size_bytes"`
	HotCacheMaxBytes  float64 `json:"hot_cache_max_bytes"`
	HotCacheEntries   int     `json:"hot_cache_entries"`
}
