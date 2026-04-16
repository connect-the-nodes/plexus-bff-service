package dto

type ConnectorsResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type FeatureResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type FeatureFlag struct {
	Name        string `json:"name" yaml:"name"`
	Parent      string `json:"parent" yaml:"parent"`
	Enabled     bool   `json:"enabled" yaml:"enabled"`
	Description string `json:"description" yaml:"description"`
}

type FeatureFlagsResponse struct {
	Features []FeatureFlag `json:"features"`
}

type FeatureList struct {
	Features []FeatureFlag `yaml:"features"`
}

type ObservabilityOverviewResponse struct {
	From              string  `json:"from,omitempty"`
	To                string  `json:"to,omitempty"`
	TotalRequests     int64   `json:"totalRequests"`
	SuccessCount      int64   `json:"successCount"`
	FailureCount      int64   `json:"failureCount"`
	SuccessRate       float64 `json:"successRate"`
	ThrottledCount    int64   `json:"throttledCount"`
	TimeoutCount      int64   `json:"timeoutCount"`
	ActiveAPICount    int     `json:"activeApiCount"`
	PeakRPS           float64 `json:"peakRps"`
	ActiveFlowCount   int     `json:"activeFlowCount"`
	IdleFlowCount     int     `json:"idleFlowCount"`
	GoldenSuccessRate float64 `json:"goldenSuccessRate"`
}

type ObservabilityInventoryItemResponse struct {
	FlowID      string  `json:"flowId"`
	APIID       string  `json:"apiId"`
	OperationID string  `json:"operationId"`
	ConnectorID string  `json:"connectorId"`
	LastSeenAt  string  `json:"lastSeenAt"`
	RequestCount int64  `json:"requestCount"`
	SuccessRate float64 `json:"successRate"`
	StatusLight string  `json:"statusLight"`
}

type AdminObservabilityOverviewResponse struct {
	From                 string  `json:"from,omitempty"`
	To                   string  `json:"to,omitempty"`
	TotalRequests        int64   `json:"totalRequests"`
	SuccessRate          float64 `json:"successRate"`
	FailureRate          float64 `json:"failureRate"`
	P95LatencyMs         float64 `json:"p95LatencyMs"`
	ThrottledCount       int64   `json:"throttledCount"`
	TimeoutCount         int64   `json:"timeoutCount"`
	ActiveTenants        int     `json:"activeTenants"`
	ActiveConnectors     int     `json:"activeConnectors"`
	TopFailingConnectorID string `json:"topFailingConnectorId"`
	TopFailingTenantID   string  `json:"topFailingTenantId"`
}

type TenantHealthSummaryResponse struct {
	TenantID         string  `json:"tenantId"`
	TotalRequests    int64   `json:"totalRequests"`
	SuccessRate      float64 `json:"successRate"`
	FailureRate      float64 `json:"failureRate"`
	QuotaUsedPercent float64 `json:"quotaUsedPercent"`
	ThrottledCount   int64   `json:"throttledCount"`
	TimeoutCount     int64   `json:"timeoutCount"`
	Status           string  `json:"status"`
}

type TenantHealthListResponse struct {
	Total int64                         `json:"total"`
	Items []TenantHealthSummaryResponse `json:"items"`
}

type ConnectorHealthSummaryResponse struct {
	ConnectorID   string  `json:"connectorId"`
	TotalRequests int64   `json:"totalRequests"`
	SuccessRate   float64 `json:"successRate"`
	FailureRate   float64 `json:"failureRate"`
	P95LatencyMs  float64 `json:"p95LatencyMs"`
	ActiveTenants int     `json:"activeTenants"`
	Status        string  `json:"status"`
}

type ConnectorHealthListResponse struct {
	Total int64                            `json:"total"`
	Items []ConnectorHealthSummaryResponse `json:"items"`
}
