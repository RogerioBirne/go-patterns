# go-patterns
The helper container with some GoF patterns for golang


## Installation

```bash
go get https://github.com/RogerioBirne/go-patterns
```

## Features
### Observer pattern
Implementation of the Observer pattern in Go.

See. https://refactoring.guru/design-patterns/observer/go/example

```go
package main

import (
	"context"
	"fmt"
	"github.com/RogerioBirne/go-patterns/gof/observer"
	"sync"
)

// eventDTO is a struct that represents the event data
type eventDTO struct {
	wg      *sync.WaitGroup
	Message string
}

// metricsObserver is a struct that represents the observer that will be notified by the subject to record some metrics
type metricsObserver struct {
	subject observer.Subject
	channel chan interface{}
}

// Notify is the method that will be called by the subject to notify the observer
func (m *metricsObserver) Notify(_ context.Context, event *eventDTO) {
	defer event.wg.Done()
	fmt.Printf("Metrics Observer: [%s]\n", event.Message)
}

// Close is the method that will be called to close the observer
func (m *metricsObserver) Close() {
	close(m.channel)
}

// newMetricsObserver is a factory method that creates a new observer
func newMetricsObserver(subject observer.Subject) observer.Observer[*eventDTO] {
	o := &metricsObserver{
		subject: subject,
		channel: make(chan interface{}, 10),
	}

	observer.StartObserver(subject, o.channel, o.Notify)
	return o
}

func main() {
	wg := &sync.WaitGroup{}

	ctx := context.TODO()
	sub := observer.NewSubject()
	o := newMetricsObserver(sub)
	defer o.Close()

	wg.Add(1)
	sub.Notify(ctx, &eventDTO{
		Message: "Observed event message",
		wg:      wg,
	})

	wg.Wait()
}
```