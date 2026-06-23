package hooksniff

import "encoding/json"

type ApiKeyResource struct{ http *httpClient }

func (r *ApiKeyResource) List() ([]ApiKey, error) {
	resp, err := r.http.request("GET", "/v1/api-keys", nil)
	if err != nil {
		return nil, err
	}
	var keys []ApiKey
	return keys, json.Unmarshal(resp, &keys)
}

func (r *ApiKeyResource) Create(name string) (*ApiKeyCreated, error) {
	resp, err := r.http.request("POST", "/v1/api-keys", map[string]string{"name": name})
	if err != nil {
		return nil, err
	}
	var key ApiKeyCreated
	return &key, json.Unmarshal(resp, &key)
}

func (r *ApiKeyResource) Delete(id string) error {
	_, err := r.http.request("DELETE", "/v1/api-keys/"+id, nil)
	return err
}

func (r *ApiKeyResource) Rotate(id string) (*ApiKeyCreated, error) {
	resp, err := r.http.request("POST", "/v1/api-keys/"+id+"/rotate", nil)
	if err != nil {
		return nil, err
	}
	var key ApiKeyCreated
	return &key, json.Unmarshal(resp, &key)
}
