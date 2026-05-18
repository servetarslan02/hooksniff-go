package models

type AppUsageStatsOut struct {
	TotalMessages      int64 `json:"totalMessages"`
	SuccessfulDeliveries int64 `json:"successfulDeliveries"`
	FailedDeliveries   int64 `json:"failedDeliveries"`
}
