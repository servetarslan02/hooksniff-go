package hooksniff

import (
	"encoding/json"
	"fmt"
)

// Paginator iterates through paginated API results.
type Paginator struct {
	http      *httpClient
	path      string
	perPage   int
	page      int
	items     []json.RawMessage
	exhausted bool
}

// newPaginator creates a new paginator for the given path.
func newPaginator(http *httpClient, path string, perPage int) *Paginator {
	if perPage <= 0 {
		perPage = 50
	}
	return &Paginator{
		http:    http,
		path:    path,
		perPage: perPage,
		page:    1,
	}
}

// Next returns the next item. Returns false when no more items.
func (p *Paginator) Next() (json.RawMessage, bool) {
	if len(p.items) == 0 && !p.exhausted {
		p.fetchPage()
	}
	if len(p.items) == 0 {
		return nil, false
	}
	item := p.items[0]
	p.items = p.items[1:]
	return item, true
}

// All collects all remaining items into a slice.
func (p *Paginator) All() ([]json.RawMessage, error) {
	var all []json.RawMessage
	for {
		item, ok := p.Next()
		if !ok {
			return all, nil
		}
		all = append(all, item)
	}
}

func (p *Paginator) fetchPage() {
	separator := "?"
	if contains(p.path, "?") {
		separator = "&"
	}
	path := fmt.Sprintf("%s%spage=%d&per_page=%d", p.path, separator, p.page, p.perPage)

	resp, err := p.http.request("GET", path, nil)
	if err != nil {
		p.exhausted = true
		return
	}

	// Try to parse as array first
	var arr []json.RawMessage
	if err := json.Unmarshal(resp, &arr); err == nil {
		p.items = arr
		if len(arr) < p.perPage {
			p.exhausted = true
		}
		p.page++
		return
	}

	// Try to parse as paginated response with "data" field
	var paginated struct {
		Data     []json.RawMessage `json:"data"`
		HasMore  *bool             `json:"has_more"`
	}
	if err := json.Unmarshal(resp, &paginated); err == nil {
		p.items = paginated.Data
		if paginated.HasMore != nil && !*paginated.HasMore {
			p.exhausted = true
		}
		if len(paginated.Data) == 0 {
			p.exhausted = true
		}
		p.page++
		return
	}

	// Try common field names
	for _, field := range []string{"deliveries", "applications", "endpoints", "templates", "schemas", "alerts"} {
		var wrapper struct {
			Items []json.RawMessage `json:"field"`
		}
		// Use a custom unmarshaler for the field
		var raw map[string]json.RawMessage
		if err := json.Unmarshal(resp, &raw); err == nil {
			if data, ok := raw[field]; ok {
				json.Unmarshal(data, &wrapper.Items)
				p.items = wrapper.Items
				if len(p.items) < p.perPage {
					p.exhausted = true
				}
				p.page++
				return
			}
		}
	}

	p.exhausted = true
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
