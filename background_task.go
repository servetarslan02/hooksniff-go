package hooksniff

import (
	"context"

	"github.com/servetarslan02/hooksniff-go/internal"
	"github.com/servetarslan02/hooksniff-go/models"
)

type BackgroundTaskListOptions struct {
	Status   *models.BackgroundTaskStatus
	Task     *models.BackgroundTaskType
	Limit    *int32
	Iterator *string
	Order    *models.Ordering
}

type BackgroundTask struct{ client *internal.HookSniffHttpClient }

func newBackgroundTask(client *internal.HookSniffHttpClient) *BackgroundTask {
	return &BackgroundTask{client: client}
}

func (b *BackgroundTask) List(ctx context.Context, options *BackgroundTaskListOptions) (*models.ListResponseBackgroundTaskOut, error) {
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
	return internal.ExecuteRequest[any, models.ListResponseBackgroundTaskOut](
		ctx, b.client, "GET", "/v1/background-task",
		nil, queryMap, nil, nil,
	)
}

func (b *BackgroundTask) Get(ctx context.Context, taskId string) (*models.BackgroundTaskOut, error) {
	pathMap := map[string]string{"task_id": taskId}
	return internal.ExecuteRequest[any, models.BackgroundTaskOut](
		ctx, b.client, "GET", "/v1/background-task/{task_id}",
		pathMap, nil, nil, nil,
	)
}
