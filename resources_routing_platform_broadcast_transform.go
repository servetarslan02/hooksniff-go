package hooksniff

import "encoding/json"

// RateLimitResource provides rate limit info.
type RateLimitResource struct{ http *httpClient }

func (r *RateLimitResource) List() (json.RawMessage, error) {
	return r.http.request("GET", "/v1/rate-limits", nil)
}

// AuditResource provides audit log access.
type AuditResource struct{ http *httpClient }

func (r *AuditResource) List() (json.RawMessage, error) {
	return r.http.request("GET", "/v1/audit-log", nil)
}

// SsoResource provides SSO configuration.
type SsoResource struct{ http *httpClient }

func (r *SsoResource) GetConfig() (json.RawMessage, error) {
	return r.http.request("GET", "/v1/sso/config", nil)
}

// CustomDomainResource manages custom domains.
type CustomDomainResource struct{ http *httpClient }

func (r *CustomDomainResource) List() (json.RawMessage, error) {
	return r.http.request("GET", "/v1/custom-domains", nil)
}

// EnvironmentResource manages environments.
type EnvironmentResource struct{ http *httpClient }

func (r *EnvironmentResource) List() (json.RawMessage, error) {
	return r.http.request("GET", "/v1/environments", nil)
}

// BroadcastResource manages broadcasts.
type BroadcastResource struct{ http *httpClient }

func (r *BroadcastResource) List() (json.RawMessage, error) {
	return r.http.request("GET", "/v1/broadcasts", nil)
}

func (r *BroadcastResource) Create(title, message string, scheduledAt ...string) (json.RawMessage, error) {
	body := map[string]string{"title": title, "message": message}
	if len(scheduledAt) > 0 {
		body["scheduled_at"] = scheduledAt[0]
	}
	return r.http.request("POST", "/v1/broadcasts", body)
}

func (r *BroadcastResource) Get(id string) (json.RawMessage, error) {
	return r.http.request("GET", "/v1/broadcasts/"+id, nil)
}

func (r *BroadcastResource) Delete(id string) error {
	_, err := r.http.request("DELETE", "/v1/broadcasts/"+id, nil)
	return err
}

func (r *BroadcastResource) Send(id string) error {
	_, err := r.http.request("POST", "/v1/broadcasts/"+id+"/send", nil)
	return err
}

// TransformResource manages payload transforms.
type TransformResource struct{ http *httpClient }

func (r *TransformResource) List(endpointID string) (json.RawMessage, error) {
	return r.http.request("GET", "/v1/endpoints/"+endpointID+"/transforms", nil)
}

func (r *TransformResource) Create(endpointID, name, code string) (json.RawMessage, error) {
	return r.http.request("POST", "/v1/endpoints/"+endpointID+"/transforms", map[string]string{"name": name, "code": code})
}

func (r *TransformResource) Get(endpointID, id string) (json.RawMessage, error) {
	return r.http.request("GET", "/v1/endpoints/"+endpointID+"/transforms/"+id, nil)
}

func (r *TransformResource) Update(endpointID, id string, params map[string]interface{}) (json.RawMessage, error) {
	return r.http.request("PUT", "/v1/endpoints/"+endpointID+"/transforms/"+id, params)
}

func (r *TransformResource) Delete(endpointID, id string) error {
	_, err := r.http.request("DELETE", "/v1/endpoints/"+endpointID+"/transforms/"+id, nil)
	return err
}
