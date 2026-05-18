package hooksniff

import (
	"context"

	"github.com/servetarslan02/hooksniff-go/internal"
	"github.com/servetarslan02/hooksniff-go/models"
)

type MessagePollerPollOptions struct {
	Limit     *int32
	Iterator  *string
	EventType *string
	Channel   *string
	After     *string
}

type MessagePoller struct{ client *internal.HookSniffHttpClient }

func newMessagePoller(client *internal.HookSniffHttpClient) *MessagePoller {
	return &MessagePoller{client: client}
}

func (m *MessagePoller) Poll(ctx context.Context, sinkId string, options *MessagePollerPollOptions) (*models.PollingEndpointOut, error) {
	pathMap := map[string]string{"sink_id": sinkId}
	queryMap := map[string]string{}
	var err error
	if options != nil {
		internal.SerializeParamToMap("limit", options.Limit, queryMap, &err)
		internal.SerializeParamToMap("iterator", options.Iterator, queryMap, &err)
		internal.SerializeParamToMap("event_type", options.EventType, queryMap, &err)
		internal.SerializeParamToMap("channel", options.Channel, queryMap, &err)
		internal.SerializeParamToMap("after", options.After, queryMap, &err)
		if err != nil {
			return nil, err
		}
	}
	return internal.ExecuteRequest[any, models.PollingEndpointOut](
		ctx, m.client, "GET", "/v1/poller/{sink_id}",
		pathMap, queryMap, nil, nil,
	)
}
