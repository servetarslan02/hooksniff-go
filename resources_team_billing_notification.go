package hooksniff

import "encoding/json"

// TeamResource manages teams and members.
type TeamResource struct{ http *httpClient }

func (r *TeamResource) List() (json.RawMessage, error) {
	return r.http.request("GET", "/v1/teams", nil)
}

func (r *TeamResource) Create(name string, description ...string) (json.RawMessage, error) {
	body := map[string]string{"name": name}
	if len(description) > 0 {
		body["description"] = description[0]
	}
	return r.http.request("POST", "/v1/teams", body)
}

func (r *TeamResource) Get(id string) (json.RawMessage, error) {
	return r.http.request("GET", "/v1/teams/"+id, nil)
}

func (r *TeamResource) Delete(id string) error {
	_, err := r.http.request("DELETE", "/v1/teams/"+id, nil)
	return err
}

func (r *TeamResource) ListMembers(teamID string) (json.RawMessage, error) {
	return r.http.request("GET", "/v1/teams/"+teamID+"/members", nil)
}

func (r *TeamResource) InviteMember(teamID, email, role string) (json.RawMessage, error) {
	return r.http.request("POST", "/v1/teams/"+teamID+"/members", map[string]string{"email": email, "role": role})
}

func (r *TeamResource) RemoveMember(teamID, memberID string) error {
	_, err := r.http.request("DELETE", "/v1/teams/"+teamID+"/members/"+memberID, nil)
	return err
}

// BillingResource manages billing and subscriptions.
type BillingResource struct{ http *httpClient }

func (r *BillingResource) Subscription() (*Subscription, error) {
	resp, err := r.http.request("GET", "/v1/billing/subscription", nil)
	if err != nil {
		return nil, err
	}
	var sub Subscription
	return &sub, json.Unmarshal(resp, &sub)
}

func (r *BillingResource) Upgrade(plan string) (json.RawMessage, error) {
	return r.http.request("POST", "/v1/billing/upgrade", map[string]string{"plan": plan})
}

func (r *BillingResource) Portal() (json.RawMessage, error) {
	return r.http.request("POST", "/v1/billing/portal", nil)
}

func (r *BillingResource) Cancel() (json.RawMessage, error) {
	return r.http.request("POST", "/v1/billing/cancel", nil)
}

func (r *BillingResource) Usage() (json.RawMessage, error) {
	return r.http.request("GET", "/v1/billing/usage", nil)
}

func (r *BillingResource) Invoices() (json.RawMessage, error) {
	return r.http.request("GET", "/v1/billing/invoices", nil)
}

// NotificationResource manages notifications.
type NotificationResource struct{ http *httpClient }

func (r *NotificationResource) List(perPage ...int) (json.RawMessage, error) {
	pp := 20
	if len(perPage) > 0 {
		pp = perPage[0]
	}
	return r.http.request("GET", "/v1/notifications?per_page="+itoa(pp), nil)
}

func (r *NotificationResource) GetUnreadCount() (int, error) {
	resp, err := r.http.request("GET", "/v1/notifications/unread-count", nil)
	if err != nil {
		return 0, err
	}
	var result struct {
		UnreadCount int `json:"unread_count"`
	}
	return result.UnreadCount, json.Unmarshal(resp, &result)
}

func (r *NotificationResource) MarkRead(id string) error {
	_, err := r.http.request("POST", "/v1/notifications/"+id+"/read", nil)
	return err
}

func (r *NotificationResource) MarkAllRead() error {
	_, err := r.http.request("POST", "/v1/notifications/read-all", nil)
	return err
}
