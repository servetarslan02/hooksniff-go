package hooksniff

import (
	"context"
	"net/http"

	"github.com/servetarslan02/hooksniff-go/internal"
	"github.com/servetarslan02/hooksniff-go/models"
)

type OperationalWebhookEndpointListOptions struct {
	Limit    *int32
	Iterator *string
	Order    *models.Ordering
}

type OperationalWebhook struct {
	Endpoint *OperationalWebhookEndpoint
}

type OperationalWebhookEndpoint struct{ client *internal.HookSniffHttpClient }

func newOperationalWebhook(client *internal.HookSniffHttpClient) *OperationalWebhook {
	return &OperationalWebhook{Endpoint: &OperationalWebhookEndpoint{client: client}}
}

func (o *OperationalWebhookEndpoint) List(ctx context.Context, options *OperationalWebhookEndpointListOptions) (*models.ListResponseOperationalWebhookEndpointOut, error) {
	queryMap := map[string]string{}
	var err error
	if options != nil {
		internal.SerializeParamToMap("limit", options.Limit, queryMap, &err)
		internal.SerializeParamToMap("iterator", options.Iterator, queryMap, &err)
		internal.SerializeParamToMap("order", options.Order, queryMap, &err)
		if err != nil {
			return nil, err
		}
	}
	return internal.ExecuteRequest[any, models.ListResponseOperationalWebhookEndpointOut](
		ctx, o.client, "GET", "/api/v1/operational-webhook/endpoint",
		nil, queryMap, nil, nil,
	)
}

func (o *OperationalWebhookEndpoint) Create(ctx context.Context, body models.OperationalWebhookEndpointIn) (*models.OperationalWebhookEndpointOut, error) {
	return internal.ExecuteRequest[models.OperationalWebhookEndpointIn, models.OperationalWebhookEndpointOut](
		ctx, o.client, "POST", "/api/v1/operational-webhook/endpoint",
		nil, nil, nil, &body,
	)
}

func (o *OperationalWebhookEndpoint) Get(ctx context.Context, endpointId string) (*models.OperationalWebhookEndpointOut, error) {
	pathMap := map[string]string{"endpoint_id": endpointId}
	return internal.ExecuteRequest[any, models.OperationalWebhookEndpointOut](
		ctx, o.client, "GET", "/api/v1/operational-webhook/endpoint/{endpoint_id}",
		pathMap, nil, nil, nil,
	)
}

func (o *OperationalWebhookEndpoint) Delete(ctx context.Context, endpointId string) error {
	pathMap := map[string]string{"endpoint_id": endpointId}
	_, err := internal.ExecuteRequest[any, any](
		ctx, o.client, http.MethodDelete, "/api/v1/operational-webhook/endpoint/{endpoint_id}",
		pathMap, nil, nil, nil,
	)
	return err
}
