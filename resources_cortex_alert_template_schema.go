package hooksniff

import "encoding/json"

// CortexResource provides Cortex AI monitoring.
type CortexResource struct{ http *httpClient }

func (r *CortexResource) Insights() (json.RawMessage, error) {
	return r.http.request("GET", "/v1/cortex/insights", nil)
}

func (r *CortexResource) Anomalies(endpointID ...string) (json.RawMessage, error) {
	path := "/v1/cortex/anomalies"
	if len(endpointID) > 0 {
		path += "?endpoint_id=" + endpointID[0]
	}
	return r.http.request("GET", path, nil)
}

func (r *CortexResource) Predict(endpointID string) (json.RawMessage, error) {
	return r.http.request("GET", "/v1/cortex/predict/"+endpointID, nil)
}

func (r *CortexResource) AutoHeal(endpointID string) (json.RawMessage, error) {
	return r.http.request("POST", "/v1/cortex/auto-heal/"+endpointID, nil)
}

// AlertResource manages alert rules.
type AlertResource struct{ http *httpClient }

func (r *AlertResource) List() (json.RawMessage, error) {
	return r.http.request("GET", "/v1/alerts", nil)
}

func (r *AlertResource) Create(name, condition string, threshold int, channels []string) (json.RawMessage, error) {
	return r.http.request("POST", "/v1/alerts", map[string]interface{}{
		"name": name, "condition": condition, "threshold": threshold, "channels": channels,
	})
}

func (r *AlertResource) Get(id string) (json.RawMessage, error) {
	return r.http.request("GET", "/v1/alerts/"+id, nil)
}

func (r *AlertResource) Update(id string, params map[string]interface{}) (json.RawMessage, error) {
	return r.http.request("PUT", "/v1/alerts/"+id, params)
}

func (r *AlertResource) Delete(id string) error {
	_, err := r.http.request("DELETE", "/v1/alerts/"+id, nil)
	return err
}

func (r *AlertResource) ListEvents(alertID string) (json.RawMessage, error) {
	return r.http.request("GET", "/v1/alerts/"+alertID+"/events", nil)
}

// TemplateResource manages payload templates.
type TemplateResource struct{ http *httpClient }

func (r *TemplateResource) List() (json.RawMessage, error) {
	return r.http.request("GET", "/v1/templates", nil)
}

func (r *TemplateResource) Create(name, content string, description ...string) (json.RawMessage, error) {
	body := map[string]string{"name": name, "content": content}
	if len(description) > 0 {
		body["description"] = description[0]
	}
	return r.http.request("POST", "/v1/templates", body)
}

func (r *TemplateResource) Get(id string) (json.RawMessage, error) {
	return r.http.request("GET", "/v1/templates/"+id, nil)
}

func (r *TemplateResource) Update(id string, params map[string]interface{}) (json.RawMessage, error) {
	return r.http.request("PUT", "/v1/templates/"+id, params)
}

func (r *TemplateResource) Delete(id string) error {
	_, err := r.http.request("DELETE", "/v1/templates/"+id, nil)
	return err
}

// SchemaResource manages the schema registry.
type SchemaResource struct{ http *httpClient }

func (r *SchemaResource) List() (json.RawMessage, error) {
	return r.http.request("GET", "/v1/schemas", nil)
}

func (r *SchemaResource) Create(name string, schema map[string]interface{}) (json.RawMessage, error) {
	return r.http.request("POST", "/v1/schemas", map[string]interface{}{"name": name, "schema": schema})
}

func (r *SchemaResource) Get(id string) (json.RawMessage, error) {
	return r.http.request("GET", "/v1/schemas/"+id, nil)
}

func (r *SchemaResource) Delete(id string) error {
	_, err := r.http.request("DELETE", "/v1/schemas/"+id, nil)
	return err
}

func (r *SchemaResource) Validate(id string, data interface{}) (json.RawMessage, error) {
	return r.http.request("POST", "/v1/schemas/"+id+"/validate", map[string]interface{}{"data": data})
}
