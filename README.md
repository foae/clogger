# clogger

Clogger is a context aware, structured logging library written in Go.

#### It has the following output formatters

1. GCP Stackdriver
2. JSON
3. Terminal (human friendly output)

#### It offers out of the box tracing

* OpenCensus
* OpenTelemetry

### How to set up a new logger

```go
package main

import log "github.com/foae/clogger"

func main() {
	// Create a new default logger
	logger := log.NewDefaultLogger()

	// We want the output to be human-readable 
	// – we choose a terminal encoder 
	logger.SetEncoder(&log.TerminalEncoder{})

	// For each and every EVENT, we can choose 
	// decorators or write our own.
	logger.SetEventOptions(log.WithOpenTelemetryTrace(), log.WithExampleEventOption())

	// For each and every LOG ENTRY, we can choose 
	// decorators or write our own.
	logger.SetLogEntryOptions(log.WithOpenCensusSpan())

	// Everything is configured, stored it globally for later use
	log.SetGlobal(logger)
}
```

### How to use the library

1. Creating new Events

```go
package main

import (
	"context"

	log "github.com/foae/clogger"
)

func handlePOST(ctx context.Context) {
	ctx, event := log.NewEvent(ctx, "Received a new POST request")
	defer event.End()

	// Attach data directly to the event
	event.Set("key", "value")
	event.SetOnErr("key", "value")
	event.SetLabel("key", "value")

	// ... TODO: add more examples, explain what each Set* is doing
}
```

2. Logging fields

```go
// consider: import log "github.com/foae/clogger"
// consider: ctx exists

log.Debug(ctx, "A debug statement")
log.With("user_id", userID).Debug(ctx, "A debug statement")
```

### Terminology

* Events
    * Set
    * SetOnErr
    * SetLabel
* Log entries
    * Fields

### Work In Progress.

Current state: `alpha`