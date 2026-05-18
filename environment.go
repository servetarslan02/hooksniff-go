package hooksniff

import (
	"context"

	"github.com/servetarslan02/hooksniff-go/internal"
	"github.com/servetarslan02/hooksniff-go/models"
)

type Environment struct{ client *internal.HookSniffHttpClient }

func newEnvironment(client *internal.HookSniffHttpClient) *Environment {
	return &Environment{client: client}
}

func (e *Environment) Export(ctx context.Context) (*models.EnvironmentOut, error) {
	return internal.ExecuteRequest[any, models.EnvironmentOut](
		ctx, e.client, "POST", "/api/v1/environment/export",
		nil, nil, nil, nil,
	)
}

func (e *Environment) Import(ctx context.Context, body models.EnvironmentIn) error {
	_, err := internal.ExecuteRequest[models.EnvironmentIn, any](
		ctx, e.client, "POST", "/api/v1/environment/import",
		nil, nil, nil, &body,
	)
	return err
}
