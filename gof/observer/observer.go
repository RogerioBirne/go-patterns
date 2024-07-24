package observer

import "context"

//go:generate mockery --name Subject --with-expecter --filename Subject.go --structname Subject
type Subject interface {
	Notify(context.Context, interface{})
	RegisterObserver(chan any)
	DeregisterObserver(chan any)
}

//go:generate mockery --name Observer --with-expecter --filename Observer.go --structname Observer
type Observer[T any] interface {
	Notify(context.Context, T)
	Close()
}
