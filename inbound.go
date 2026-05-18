package hooksniff

import (
	"context"
	"net/http"

	"github.com/servetarslan02/hooksniff-go/internal"
)

type InboundConfig struct {
	Id         string  `json:"id"`
	CustomerId string  `json:"customer_id"`
	Provider   string  `json:"provider"`
	Secret     string  `json:"secret"`
	EndpointId *string `json:"endpoint_id,omitempty"`
	Enabled    bool    `json:"enabled"`
	CreatedAt  string  `json:"created_at"`
	UpdatedAt  string  `json:"updated_at"`
}

type InboundConfigIn struct {
	Provider   string  `json:"provider"`
	Secret     string  `json:"secret"`
	EndpointId *string `json:"endpoint_id,omitempty"`
	Enabled    *bool   `json:"enabled,omitempty"`
}

type InboundListOptions struct {
	Limit    *int32
	Iterator *string
}

type Inbound struct{ client *internal.HookSniffHttpClient }

func newInbound(client *internal.HookSniffHttpClient) *Inbound {
	return &Inbound{client: client}
}

func (i *Inbound) List(ctx context.Context, appId string, options *InboundListOptions) ([]InboundConfig, error) {
	pathMap := map[string]string{"app_id": appId}
	queryMap := map[string]string{}
	var err error
	if options != nil {
		internal.SerializeParamToMap("limit", options.Limit, queryMap, &err)
		internal.SerializeParamToMap("iterator", options.Iterator, queryMap, &err)
		if err != nil {
			return nil, err
		}
	}
	result, err := internal.ExecuteRequest[any, []InboundConfig](
		ctx, i.client, "GET", "/v1/app/{app_id}/inbound",
		pathMap, queryMap, nil, nil,
	)
	if err != nil {
		return nil, err
	}
	return *result, nil
}

func (i *Inbound) Create(ctx context.Context, appId string, body InboundConfigIn) (*InboundConfig, error) {
	pathMap := map[string]string{"app_id": appId}
	return internal.ExecuteRequest[InboundConfigIn, InboundConfig](
		ctx, i.client, "POST", "/v1/app/{app_id}/inbound",
		pathMap, nil, nil, &body,
	)
}

func (i *Inbound) Get(ctx context.Context, appId, inboundId string) (*InboundConfig, error) {
	pathMap := map[string]string{"app_id": appId, "inbound_id": inboundId}
	return internal.ExecuteRequest[any, InboundConfig](
		ctx, i.client, "GET", "/v1/app/{app_id}/inbound/{inbound_id}",
		pathMap, nil, nil, nil,
	)
}

func (i *Inbound) Delete(ctx context.Context, appId, inboundId string) error {
	pathMap := map[string]string{"app_id": appId, "inbound_id": inboundId}
	_, err := internal.ExecuteRequest[any, any](
		ctx, i.client, http.MethodDelete, "/v1/app/{app_id}/inbound/{inbound_id}",
		pathMap, nil, nil, nil,
	)
	return err
}
