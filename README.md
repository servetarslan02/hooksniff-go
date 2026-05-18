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

    endpoints, err := client.Endpoint.List(context.Background(), nil)
    if err != nil {
        panic(err)
    }
    for _, ep := range endpoints.Data {
        fmt.Printf("Endpoint: %s\n", ep.Id)
    }
}
```

## API Resources

| Resource | Methods |
|----------|---------|
| `Endpoint` | List, Create, Get, Update, Delete, Patch |
| `Message` | Create, Get, List, Expunge |
| `MessageAttempt` | List, Get, Resend, ListByMsg, ListByEndpoint |
| `EventType` | List, Create, Get, Update, Delete |
| `Authentication` | Logout |
| `Environment` | Export, Import |
| `BackgroundTask` | List, Get |
| `Connector` | List, Get, Delete |
| `Integration` | List, Get, Create, Update, Delete, RotateKey |
| `Inbound` | List, Create, Get, Delete |
| `Stream` | ListChannels, GetChannel, CreateChannel, Subscribe, Publish |
| `MessagePoller` | Poll |
| `OperationalWebhook` | Endpoint.List, Endpoint.Create, Endpoint.Get, Endpoint.Delete |
| `Statistics` | AggregateEventTypes |

## Webhook Verification

```go
wh, err := hooksniff.NewWebhook("whsec_...")
if err != nil {
    panic(err)
}
err = wh.Verify(payload, headers)
```

## Features

- ✅ HMAC-SHA256 webhook verification
- ✅ Typed webhook events
- ✅ Automatic retry with exponential backoff
- ✅ Pagination helpers
- ✅ Rate limit header parsing
- ✅ SSE streaming
- ✅ Idempotency keys
- ✅ Configurable HTTP client

## License

MIT
