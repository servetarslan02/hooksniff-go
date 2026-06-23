package hooksniff

import "encoding/json"

// BackgroundTaskResource manages background tasks.
type BackgroundTaskResource struct{ http *httpClient }

func (r *BackgroundTaskResource) List() (json.RawMessage, error) {
	return r.http.request("GET", "/v1/background-tasks", nil)
}

func (r *BackgroundTaskResource) Get(id string) (json.RawMessage, error) {
	return r.http.request("GET", "/v1/background-tasks/"+id, nil)
}

func (r *BackgroundTaskResource) Cancel(id string) error {
	_, err := r.http.request("POST", "/v1/background-tasks/"+id+"/cancel", nil)
	return err
}

// IntegrationResource manages integrations.
type IntegrationResource struct{ http *httpClient }

func (r *IntegrationResource) List() (json.RawMessage, error) {
	return r.http.request("GET", "/v1/integrations", nil)
}

func (r *IntegrationResource) Get(id string) (json.RawMessage, error) {
	return r.http.request("GET", "/v1/integrations/"+id, nil)
}

func (r *IntegrationResource) Delete(id string) error {
	_, err := r.http.request("DELETE", "/v1/integrations/"+id, nil)
	return err
}

func (r *IntegrationResource) RotateKey(id string) (json.RawMessage, error) {
	return r.http.request("POST", "/v1/integrations/"+id+"/rotate-key", nil)
}

// ServiceTokenResource manages service tokens.
type ServiceTokenResource struct{ http *httpClient }

func (r *ServiceTokenResource) List() (json.RawMessage, error) {
	return r.http.request("GET", "/v1/service-tokens", nil)
}

func (r *ServiceTokenResource) Create(name string) (json.RawMessage, error) {
	return r.http.request("POST", "/v1/service-tokens", map[string]string{"name": name})
}

func (r *ServiceTokenResource) Delete(id string) error {
	_, err := r.http.request("DELETE", "/v1/service-tokens/"+id, nil)
	return err
}

// OperationalWebhookResource manages operational webhooks.
type OperationalWebhookResource struct{ http *httpClient }

func (r *OperationalWebhookResource) List() (json.RawMessage, error) {
	return r.http.request("GET", "/v1/operational-webhooks", nil)
}

func (r *OperationalWebhookResource) Create(url string, events []string) (json.RawMessage, error) {
	return r.http.request("POST", "/v1/operational-webhooks", map[string]interface{}{"url": url, "events": events})
}

func (r *OperationalWebhookResource) Get(id string) (json.RawMessage, error) {
	return r.http.request("GET", "/v1/operational-webhooks/"+id, nil)
}

func (r *OperationalWebhookResource) Delete(id string) error {
	_, err := r.http.request("DELETE", "/v1/operational-webhooks/"+id, nil)
	return err
}
