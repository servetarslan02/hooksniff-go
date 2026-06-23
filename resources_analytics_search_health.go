package hooksniff

import (
	"encoding/json"
	"strconv"
)

// AnalyticsResource provides delivery analytics.
type AnalyticsResource struct{ http *httpClient }

func (r *AnalyticsResource) Deliveries(rangeParam ...string) (json.RawMessage, error) {
	rng := "24h"
	if len(rangeParam) > 0 {
		rng = rangeParam[0]
	}
	return r.http.request("GET", "/v1/analytics/deliveries?range="+rng, nil)
}

func (r *AnalyticsResource) SuccessRate(rangeParam ...string) (json.RawMessage, error) {
	rng := "24h"
	if len(rangeParam) > 0 {
		rng = rangeParam[0]
	}
	return r.http.request("GET", "/v1/analytics/success-rate?range="+rng, nil)
}

func (r *AnalyticsResource) Latency(rangeParam ...string) (json.RawMessage, error) {
	rng := "24h"
	if len(rangeParam) > 0 {
		rng = rangeParam[0]
	}
	return r.http.request("GET", "/v1/analytics/latency?range="+rng, nil)
}

// SearchResource provides delivery search.
type SearchResource struct{ http *httpClient }

func (r *SearchResource) Deliveries(query string, page, perPage int) (*SearchResult, error) {
	if page <= 0 { page = 1 }
	if perPage <= 0 { perPage = 20 }
	resp, err := r.http.request("GET", "/v1/search?q="+query+"&page="+itoa(page)+"&per_page="+itoa(perPage), nil)
	if err != nil {
		return nil, err
	}
	var result SearchResult
	return &result, json.Unmarshal(resp, &result)
}

// HealthResource provides health checks and outbound IPs.
type HealthResource struct{ http *httpClient }

func (r *HealthResource) Check() (*HealthResponse, error) {
	resp, err := r.http.request("GET", "/health", nil)
	if err != nil {
		return nil, err
	}
	var h HealthResponse
	return &h, json.Unmarshal(resp, &h)
}

func (r *HealthResource) OutboundIPs() (*OutboundIPs, error) {
	resp, err := r.http.request("GET", "/v1/outbound-ips", nil)
	if err != nil {
		return nil, err
	}
	var ips OutboundIPs
	return &ips, json.Unmarshal(resp, &ips)
}

func itoa(i int) string {
	return strconv.Itoa(i)
}
