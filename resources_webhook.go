package hooksniff

import "encoding/json"

type WebhookResource struct{ http *httpClient }

func (r *WebhookResource) Send(params *WebhookSend, opts ...requestOption) (*WebhookDelivery, error) {
	resp, err := r.http.request("POST", "/v1/webhooks", params, opts...)
	if err != nil {
		return nil, err
	}
	var d WebhookDelivery
	return &d, json.Unmarshal(resp, &d)
}

func (r *WebhookResource) SendBatch(webhooks []WebhookSend, opts ...requestOption) ([]WebhookDelivery, error) {
	resp, err := r.http.request("POST", "/v1/webhooks/batch", map[string]interface{}{"webhooks": webhooks}, opts...)
	if err != nil {
		return nil, err
	}
	var result struct {
		Deliveries []WebhookDelivery `json:"deliveries"`
	}
	return result.Deliveries, json.Unmarshal(resp, &result)
}

func (r *WebhookResource) List(perPage ...int) *Paginator {
	pp := 50
	if len(perPage) > 0 {
		pp = perPage[0]
	}
	return newPaginator(r.http, "/v1/webhooks", pp)
}

func (r *WebhookResource) Get(id string) (*WebhookDelivery, error) {
	resp, err := r.http.request("GET", "/v1/webhooks/"+id, nil)
	if err != nil {
		return nil, err
	}
	var d WebhookDelivery
	return &d, json.Unmarshal(resp, &d)
}

func (r *WebhookResource) Replay(id string) (*WebhookDelivery, error) {
	resp, err := r.http.request("POST", "/v1/webhooks/"+id+"/replay", nil)
	if err != nil {
		return nil, err
	}
	var d WebhookDelivery
	return &d, json.Unmarshal(resp, &d)
}

func (r *WebhookResource) BatchReplay(ids []string) (int, error) {
	resp, err := r.http.request("POST", "/v1/webhooks/batch-replay", map[string]interface{}{"webhook_ids": ids})
	if err != nil {
		return 0, err
	}
	var result struct {
		Replayed int `json:"replayed"`
	}
	return result.Replayed, json.Unmarshal(resp, &result)
}
