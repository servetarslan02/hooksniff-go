package hooksniff

import "encoding/json"

// ConnectorResource manages connectors and their configs.
type ConnectorResource struct{ http *httpClient }

func (r *ConnectorResource) List() (json.RawMessage, error) {
	return r.http.request("GET", "/v1/connectors", nil)
}

func (r *ConnectorResource) Get(id string) (json.RawMessage, error) {
	return r.http.request("GET", "/v1/connectors/"+id, nil)
}

func (r *ConnectorResource) ListConfigs() (json.RawMessage, error) {
	return r.http.request("GET", "/v1/connectors/configs", nil)
}

func (r *ConnectorResource) CreateConfig(connectorID, name string, config map[string]interface{}) (json.RawMessage, error) {
	return r.http.request("POST", "/v1/connectors/configs", map[string]interface{}{
		"connector_id": connectorID, "name": name, "config": config,
	})
}

func (r *ConnectorResource) GetConfig(id string) (json.RawMessage, error) {
	return r.http.request("GET", "/v1/connectors/configs/"+id, nil)
}

func (r *ConnectorResource) UpdateConfig(id string, params map[string]interface{}) (json.RawMessage, error) {
	return r.http.request("PUT", "/v1/connectors/configs/"+id, params)
}

func (r *ConnectorResource) DeleteConfig(id string) error {
	_, err := r.http.request("DELETE", "/v1/connectors/configs/"+id, nil)
	return err
}

// StreamResource manages real-time streaming channels.
type StreamResource struct{ http *httpClient }

func (r *StreamResource) ListChannels() (json.RawMessage, error) {
	return r.http.request("GET", "/v1/stream/channels", nil)
}

func (r *StreamResource) CreateChannel(name string, description ...string) (json.RawMessage, error) {
	body := map[string]string{"name": name}
	if len(description) > 0 {
		body["description"] = description[0]
	}
	return r.http.request("POST", "/v1/stream/channels", body)
}

func (r *StreamResource) GetChannel(id string) (json.RawMessage, error) {
	return r.http.request("GET", "/v1/stream/channels/"+id, nil)
}

func (r *StreamResource) DeleteChannel(id string) error {
	_, err := r.http.request("DELETE", "/v1/stream/channels/"+id, nil)
	return err
}

func (r *StreamResource) ListMessages(channelID string) (json.RawMessage, error) {
	return r.http.request("GET", "/v1/stream/channels/"+channelID+"/messages", nil)
}

func (r *StreamResource) Publish(channelID, event string, data map[string]interface{}) (json.RawMessage, error) {
	return r.http.request("POST", "/v1/stream/publish", map[string]interface{}{
		"channel_id": channelID, "event": event, "data": data,
	})
}

func (r *StreamResource) ListSubscriptions() (json.RawMessage, error) {
	return r.http.request("GET", "/v1/stream/subscriptions", nil)
}

func (r *StreamResource) GetSubscription(id string) (json.RawMessage, error) {
	return r.http.request("GET", "/v1/stream/subscriptions/"+id, nil)
}

func (r *StreamResource) DisconnectSubscription(id string) error {
	_, err := r.http.request("DELETE", "/v1/stream/subscriptions/"+id, nil)
	return err
}
