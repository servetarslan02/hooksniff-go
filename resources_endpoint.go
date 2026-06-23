package hooksniff

import "encoding/json"

type EndpointResource struct{ http *httpClient }

func (r *EndpointResource) Create(params *EndpointCreate) (*Endpoint, error) {
	resp, err := r.http.request("POST", "/v1/endpoints", params)
	if err != nil {
		return nil, err
	}
	var ep Endpoint
	return &ep, json.Unmarshal(resp, &ep)
}

func (r *EndpointResource) List(perPage ...int) *Paginator {
	pp := 50
	if len(perPage) > 0 {
		pp = perPage[0]
	}
	return newPaginator(r.http, "/v1/endpoints", pp)
}

func (r *EndpointResource) Get(id string) (*Endpoint, error) {
	resp, err := r.http.request("GET", "/v1/endpoints/"+id, nil)
	if err != nil {
		return nil, err
	}
	var ep Endpoint
	return &ep, json.Unmarshal(resp, &ep)
}

func (r *EndpointResource) Update(id string, params *EndpointUpdate) (*Endpoint, error) {
	resp, err := r.http.request("PUT", "/v1/endpoints/"+id, params)
	if err != nil {
		return nil, err
	}
	var ep Endpoint
	return &ep, json.Unmarshal(resp, &ep)
}

func (r *EndpointResource) Delete(id string) error {
	_, err := r.http.request("DELETE", "/v1/endpoints/"+id, nil)
	return err
}

func (r *EndpointResource) RotateSecret(id string) (*SecretRotateResponse, error) {
	resp, err := r.http.request("POST", "/v1/endpoints/"+id+"/rotate-secret", nil)
	if err != nil {
		return nil, err
	}
	var sr SecretRotateResponse
	return &sr, json.Unmarshal(resp, &sr)
}
