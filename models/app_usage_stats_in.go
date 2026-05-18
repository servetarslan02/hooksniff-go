package models

type AppUsageStatsIn struct {
	Since *string `json:"since,omitempty"`
	Until *string `json:"until,omitempty"`
}
