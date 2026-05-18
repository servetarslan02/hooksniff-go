package hooksniff

import (
	"context"

	"github.com/servetarslan02/hooksniff-go/models"
)

// ListAll returns a Paginator that auto-paginates through all messages.
//
// Usage:
//
//	paginator := hs.Message.ListAll(ctx, appId, &hooksniff.MessageListOptions{Limit: lo(100)})
//	for paginator.Next() {
//	    msg := paginator.Value()
//	    fmt.Println(msg.Id)
//	}
//	if err := paginator.Err(); err != nil {
//	    log.Fatal(err)
//	}
func (message *Message) ListAll(ctx context.Context, appId string, o *MessageListOptions) *Paginator[models.MessageOut] {
	return NewPaginator(ctx, func(ctx context.Context, iterator *string) (ListResponse[models.MessageOut], error) {
		opts := &MessageListOptions{}
		if o != nil {
			*opts = *o
		}
		opts.Iterator = iterator
		return message.List(ctx, appId, opts)
	})
}

// Instantiates a new MessageIn object with a raw string payload.
// The payload is not normalized on the server. Normally, payloads are required
// to be JSON, and HookSniff will minify the payload before sending the webhook
// (for example, by removing extraneous whitespace or unnecessarily escaped
// characters in strings). With this function, the payload will be sent
// "as is", without any minification or other processing.
//
// The `contentType` parameter can be used to change the `content-type` header
// of the webhook sent by HookSniff overriding the default of `application/json`.
//
// See the class documentation for details about the other parameters.
func NewMessageInRaw(
	eventType string,
	payload string,
	contentType *string,
) *models.MessageIn {
	msgIn := models.MessageIn{
		EventType: eventType,
		Payload:   map[string]any{},
	}

	transformationsParams := map[string]interface{}{
		"rawPayload": payload,
	}
	if contentType != nil {
		transformationsParams["headers"] = map[string]string{
			"content-type": *contentType,
		}
	}
	msgIn.TransformationsParams = &transformationsParams

	return &msgIn
}
