// Package hooksniff this file is @generated DO NOT EDIT
package hooksniff

import (
	"context"

	"github.com/servetarslan02/hooksniff-go/internal"
	"github.com/servetarslan02/hooksniff-go/models"
)

type Statistics struct {
	client *internal.HookSniffHttpClient
}

func newStatistics(client *internal.HookSniffHttpClient) *Statistics {
	return &Statistics{
		client: client,
	}
}

// Creates a background task to calculate the listed event types for all apps in the organization.
//
// Note that this endpoint is asynchronous. You will need to poll the `Get Background Task` endpoint to
// retrieve the results of the operation.
//
// The completed background task will return a payload like the following:
// ```json
//
//	{
//	  "id": "qtask_33qe39Stble9Rn3ZxFrqL5ZSsjT",
//	  "status": "finished",
//	  "task": "event-type.aggregate",
//	  "data": {
//	    "event_types": [
//	      {
//	        "appId": "app_33W1An2Zz5cO9SWbhHsYyDmVC6m",
//	        "explicitlySubscribedEventTypes": ["user.signup", "user.deleted"],
//	        "hasCatchAllEndpoint": false
//	      }
//	    ]
//	  }
//	}
//
// ```
func (statistics *Statistics) AggregateEventTypes(
	ctx context.Context,
) (*models.AggregateEventTypesOut, error) {
	return internal.ExecuteRequest[any, models.AggregateEventTypesOut](
		ctx,
		statistics.client,
		"PUT",
		"/v1/stats/usage/event-types",
		nil,
		nil,
		nil,
		nil,
	)
}
