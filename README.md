# HookSniff Go SDK

<p align="center">
  <a href="https://pkg.go.dev/github.com/servetarslan02/hooksniff-go"><img src="https://pkg.go.dev/badge/github.com/servetarslan02/hooksniff-go.svg" alt="Go Reference"></a>
  <a href="https://github.com/servetarslan02/HookSniff"><img src="https://img.shields.io/github/license/servetarslan02/HookSniff" alt="License"></a>
</p>

Go SDK for the [HookSniff](https://hooksniff.vercel.app) webhook delivery platform.

## Installation

```bash
go get github.com/servetarslan02/hooksniff-go@v1.2.0
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    hooksniff "github.com/servetarslan02/hooksniff-go"
)

func main() {
    client, err := hooksniff.New("your-api-token", nil)
    if err != nil {
        panic(err)
    }

    // List applications
    apps, err := client.Endpoint.List(context.Background(), "app_id", nil)
    if err != nil {
        panic(err)
    }
    for _, app := range apps.Data {
        fmt.Printf("Endpoint: %s\n", app.Id)
    }
}
```

## API Resources

| Resource | Methods | Description |
|----------|---------|-------------|
| `Application` | List, Create, Get, Update, Delete, Patch | Manage applications |
| `Endpoint` | List, Create, Get, Update, Delete, Patch | Manage webhook endpoints |
| `Message` | Create, Get, List, Expunge | Send and manage messages |
| `MessageAttempt` | List, Get, Resend, ListByMsg, ListByEndpoint | Track delivery attempts |
| `EventType` | List, Create, Get, Update, Delete | Manage event types |
| `Authentication` | AppPortalAccess, Logout, ExpireAll | Authentication management |
| `Environment` | Export, Import | Organization settings |
| `BackgroundTask` | List, Get | Background task status |
| `Connector` | List, Get, Delete | Third-party connectors |
| `Integration` | List, Get, Create, Update, Delete, RotateKey | App integrations |
| `Inbound` | List, Create, Get, Delete | Inbound webhook configs |
| `Stream` | ListChannels, GetChannel, CreateChannel, UpdateChannel, DeleteChannel, Publish | Real-time streaming |
| `MessagePoller` | Poll | Poll-based message consumption |
| `OperationalWebhook` | Endpoint.List, Endpoint.Create, Endpoint.Get, Endpoint.Delete | Operational webhook endpoints |
| `Statistics` | AggregateAppStats | Usage statistics |

## Webhook Verification

```go
import "github.com/servetarslan02/hooksniff-go"

wh, err := hooksniff.NewWebhook("whsec_...")
if err != nil {
    panic(err)
}

// Verify webhook signature
err = wh.Verify(payload, headers)
if err != nil {
    // Invalid signature
}
```

## License

MIT
