package hooksniff

import (
	"context"
	"net/http"

	"github.com/servetarslan02/hooksniff-go/internal"
	"github.com/servetarslan02/hooksniff-go/models"
)

type ConnectorListOptions struct {
	Limit       *int32
	Iterator    *string
	Order       *models.Ordering
	ProductType *models.ConnectorProduct
}

type ConnectorApi struct{ client *internal.HookSniffHttpClient }

func newConnectorApi(client *internal.HookSniffHttpClient) *ConnectorApi {
	return &ConnectorApi{client: client}
}

func (c *ConnectorApi) List(ctx context.Context, options *ConnectorListOptions) (*models.ListResponseConnectorOut, error) {
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
	return internal.ExecuteRequest[any, models.ListResponseConnectorOut](
		ctx, c.client, "GET", "/api/v1/connector",
		nil, queryMap, nil, nil,
	)
}

func (c *ConnectorApi) Get(ctx context.Context, connectorId string) (*models.ConnectorOut, error) {
	pathMap := map[string]string{"connector_id": connectorId}
	return internal.ExecuteRequest[any, models.ConnectorOut](
		ctx, c.client, "GET", "/api/v1/connector/{connector_id}",
		pathMap, nil, nil, nil,
	)
}

func (c *ConnectorApi) Delete(ctx context.Context, connectorId string) error {
	pathMap := map[string]string{"connector_id": connectorId}
	_, err := internal.ExecuteRequest[any, any](
		ctx, c.client, http.MethodDelete, "/api/v1/connector/{connector_id}",
		pathMap, nil, nil, nil,
	)
	return err
}
