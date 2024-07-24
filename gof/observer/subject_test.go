package observer_test

import (
	"context"
	"fmt"
	"github.com/RogerioBirne/go-patterns/gof/observer"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

// eventDTO is a struct that represents the event data.
type eventDTO struct {
	t       *testing.T
	wg      *sync.WaitGroup
	Message string
}

// metricsObserver is a struct that represents the observer that will be notified by the subject to record some metrics.
type metricsObserver struct {
	subject observer.Subject
	channel chan interface{}
}

// Notify is the method that will be called by the subject to notify the observer.
func (m *metricsObserver) Notify(_ context.Context, event *eventDTO) {
	defer event.wg.Done()
	assert.Equal(event.t, "Observed event message", event.Message)
	fmt.Printf("Metrics Observer: [%s]\n", event.Message)
}

// Close is the method that will be called to close the observer.
func (m *metricsObserver) Close() {
	close(m.channel)
}

// newMetricsObserver is a factory method that creates a new observer.
func newMetricsObserver(subject observer.Subject) observer.Observer[*eventDTO] {
	o := &metricsObserver{
		subject: subject,
		channel: make(chan interface{}, 10),
	}

	observer.StartObserver(subject, o.channel, o.Notify)
	return o
}

func TestObserverPattern(t *testing.T) {
	wg := &sync.WaitGroup{}

	ctx := context.TODO()
	sub := observer.NewSubject()
	o := newMetricsObserver(sub)
	defer o.Close()

	wg.Add(1)
	sub.Notify(ctx, &eventDTO{
		Message: "Observed event message",
		t:       t,
		wg:      wg,
	})

	sub.Notify(ctx, "Not eventDTO instance")

	wg.Wait()
}
