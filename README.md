# HookSniff Go SDK

Official Go SDK for [HookSniff](https://hooksniff.vercel.app) — the webhook infrastructure for developers.

## Installation

```bash
go get github.com/servetarslan02/hooksniff-go
```

## Quick Start

```go
package main

import (
    "fmt"
    hooksniff "github.com/servetarslan02/hooksniff-go"
)

func main() {
    hs := hooksniff.New("hr_live_...", nil)

    // Create an application
    app, _ := hs.Application.Create(&hooksniff.ApplicationCreate{
        Name: "My App",
    })

    // Create an endpoint
    ep, _ := hs.Endpoint.Create(&hooksniff.EndpointCreate{
        URL:           "https://app.com/webhook",
        ApplicationID: app.ID,
    })

    // Send a webhook
    delivery, _ := hs.Webhook.Send(&hooksniff.WebhookSend{
        EndpointID: ep.ID,
        Event:      "order.created",
        Data:       map[string]interface{}{"order_id": "12345"},
    })

    fmt.Println(delivery.ID)
}
```

## Features

- **Zero dependencies** — uses Go standard library only
- **Full type safety** — all structs are typed
- **Auto-retry** — exponential backoff on 429/5xx errors
- **Auto-pagination** — iterate through all resources
- **Webhook verification** — Standard Webhooks compliant
- **27 resources** — Application, Endpoint, Webhook, Billing, Cortex, Teams, and more

## Usage

### Webhook Verification

```go
wh, _ := hooksniff.NewWebhook("whsec_...")

event, err := wh.Verify(payload, headers)
if err != nil {
    // Invalid signature
}
```

### Error Handling

```go
delivery, err := hs.Webhook.Send(params)
if err != nil {
    switch e := err.(type) {
    case *hooksniff.AuthenticationError:
        fmt.Println("Invalid API key")
    case *hooksniff.NotFoundError:
        fmt.Println("Resource not found")
    case *hooksniff.RateLimitError:
        fmt.Printf("Rate limited, retry after %ds\n", e.RetryAfter)
    case *hooksniff.ValidationError:
        fmt.Println("Validation error:", e.Detail)
    }
}
```

## License

MIT — see [LICENSE](LICENSE) for details.
