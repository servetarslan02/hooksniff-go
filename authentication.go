// Package hooksniff this file is @generated DO NOT EDIT
package hooksniff

import (
	"context"

	"github.com/servetarslan02/hooksniff-go/internal"
	"github.com/servetarslan02/hooksniff-go/models"
)

type Authentication struct {
	client *internal.HookSniffHttpClient
}

func newAuthentication(client *internal.HookSniffHttpClient) *Authentication {
	return &Authentication{
		client: client,
	}
}

type AuthenticationLogoutOptions struct {
	IdempotencyKey *string
}

type AuthenticationStreamLogoutOptions struct {
	IdempotencyKey *string
}

type AuthenticationRotateStreamPollerTokenOptions struct {
	IdempotencyKey *string
}

// Logout an app token.
//
// Trying to log out other tokens will fail.
func (authentication *Authentication) Logout(
	ctx context.Context,
	o *AuthenticationLogoutOptions,
) error {
	headerMap := map[string]string{}
	var err error
	if o != nil {
		internal.SerializeParamToMap("idempotency-key", o.IdempotencyKey, headerMap, &err)
		if err != nil {
			return err
		}
	}
	_, err = internal.ExecuteRequest[any, any](
		ctx,
		authentication.client,
		"POST",
		"/v1/auth/logout",
		nil,
		nil,
		headerMap,
		nil,
	)
	return err
}

// Logout a stream token.
//
// Trying to log out other tokens will fail.
func (authentication *Authentication) StreamLogout(
	ctx context.Context,
	o *AuthenticationStreamLogoutOptions,
) error {
	headerMap := map[string]string{}
	var err error
	if o != nil {
		internal.SerializeParamToMap("idempotency-key", o.IdempotencyKey, headerMap, &err)
		if err != nil {
			return err
		}
	}
	_, err = internal.ExecuteRequest[any, any](
		ctx,
		authentication.client,
		"POST",
		"/v1/auth/stream-logout",
		nil,
		nil,
		headerMap,
		nil,
	)
	return err
}

// Get the current auth token for the stream poller.
func (authentication *Authentication) GetStreamPollerToken(
	ctx context.Context,
	streamId string,
	sinkId string,
) (*models.ApiTokenOut, error) {
	pathMap := map[string]string{
		"stream_id": streamId,
		"sink_id":   sinkId,
	}
	return internal.ExecuteRequest[any, models.ApiTokenOut](
		ctx,
		authentication.client,
		"GET",
		"/v1/auth/stream/{stream_id}/sink/{sink_id}/poller/token",
		pathMap,
		nil,
		nil,
		nil,
	)
}

// Create a new auth token for the stream poller API.
func (authentication *Authentication) RotateStreamPollerToken(
	ctx context.Context,
	streamId string,
	sinkId string,
	rotatePollerTokenIn models.RotatePollerTokenIn,
	o *AuthenticationRotateStreamPollerTokenOptions,
) (*models.ApiTokenOut, error) {
	pathMap := map[string]string{
		"stream_id": streamId,
		"sink_id":   sinkId,
	}
	headerMap := map[string]string{}
	var err error
	if o != nil {
		internal.SerializeParamToMap("idempotency-key", o.IdempotencyKey, headerMap, &err)
		if err != nil {
			return nil, err
		}
	}
	return internal.ExecuteRequest[models.RotatePollerTokenIn, models.ApiTokenOut](
		ctx,
		authentication.client,
		"POST",
		"/v1/auth/stream/{stream_id}/sink/{sink_id}/poller/token/rotate",
		pathMap,
		nil,
		headerMap,
		&rotatePollerTokenIn,
	)
}
