package observer

import "context"

type subject struct {
	observers []chan interface{}
}

type data struct {
	ctx   context.Context
	event any
}

func NewSubject() Subject {
	return &subject{
		observers: make([]chan interface{}, 0),
	}
}

func (s *subject) Notify(ctx context.Context, e interface{}) {
	for _, o := range s.observers {
		o <- data{
			ctx:   ctx,
			event: e,
		}
	}
}

func (s *subject) RegisterObserver(observerCh chan interface{}) {
	s.observers = append(s.observers, observerCh)
}

func (s *subject) DeregisterObserver(observerCh chan interface{}) {
	for i, ch := range s.observers {
		if ch == observerCh {
			s.observers = append(s.observers[:i], s.observers[i+1:]...)
			break
		}
	}
}

func StartObserver[T any](s Subject, ch chan interface{}, fnNotify func(context.Context, T)) {
	s.RegisterObserver(ch)

	go readLoopCh(s, ch, fnNotify)
}

func readLoopCh[T any](subject Subject, channel chan interface{}, notify func(context.Context, T)) {
	for {
		if msg, ok := <-channel; !ok {
			subject.DeregisterObserver(channel)
			return
		} else if value, ok := msg.(data); ok {
			if event, isT := value.event.(T); isT {
				notify(value.ctx, event)
			}
		}
	}
}
