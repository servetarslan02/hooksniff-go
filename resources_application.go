package hooksniff

import "encoding/json"

// ApplicationResource manages HookSniff applications.
type ApplicationResource struct {
	http *httpClient
}

func (r *ApplicationResource) Create(params *ApplicationCreate) (*Application, error) {
	resp, err := r.http.request("POST", "/v1/applications", params)
	if err != nil {
		return nil, err
	}
	var app Application
	return &app, json.Unmarshal(resp, &app)
}

func (r *ApplicationResource) List(perPage ...int) *Paginator {
	pp := 50
	if len(perPage) > 0 {
		pp = perPage[0]
	}
	return newPaginator(r.http, "/v1/applications", pp)
}

func (r *ApplicationResource) Get(id string) (*Application, error) {
	resp, err := r.http.request("GET", "/v1/applications/"+id, nil)
	if err != nil {
		return nil, err
	}
	var app Application
	return &app, json.Unmarshal(resp, &app)
}

func (r *ApplicationResource) Update(id string, params *ApplicationUpdate) (*Application, error) {
	resp, err := r.http.request("PUT", "/v1/applications/"+id, params)
	if err != nil {
		return nil, err
	}
	var app Application
	return &app, json.Unmarshal(resp, &app)
}

func (r *ApplicationResource) Delete(id string) error {
	_, err := r.http.request("DELETE", "/v1/applications/"+id, nil)
	return err
}
