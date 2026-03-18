package ungo

import (
	"sync/atomic"
	"time"
)

type Event struct {
	Type    int
	Data    any
	ReplyTo Optional[chan any]
}

type EventLoop struct {
	events          chan Event
	stop            chan struct{}
	handlers        FastMap[int, func(any) error]
	handleUnhandled func(Event)
	running         uint32
}

func NewEventLoop() *EventLoop {
	return &EventLoop{
		events:          make(chan Event),
		stop:            make(chan struct{}),
		handlers:        FastMap[int, func(any) error]{},
		handleUnhandled: func(ev Event) {},
		running:         0,
	}
}

func (el *EventLoop) ToSubProcess() SubProcess {
	return SubProcess{
		fn: func() Optional[int] {
			el.Run()
			return MakeOptional(int(ExitSuccess))
		},
		on_finish: Optional[func(Optional[int])]{},
	}
}

type Future[T any] struct {
	result Optional[T]
	done   chan struct{}
}

func (f *Future[T]) Get() Optional[T] {
	<-f.done
	return f.result
}

func (f *Future[T]) Set(val T) {
	f.result = MakeOptional(val)
	close(f.done)
}

func (el *EventLoop) Stop() {
	close(el.stop)
}

func (el *EventLoop) Start() {
	if !atomic.CompareAndSwapUint32(&el.running, 0, 1) {
		return
	}

	go func() {
		OnShutdown(func() {
			el.Stop()
		})

		for {
			select {
			case ev := <-el.events:
				el.process(ev)
			case <-el.stop:
				return
			}
		}
	}()
}

func (el *EventLoop) process(ev Event) {
	handler, ok := el.handlers.Get(ev.Type)
	if !ok {
		el.handleUnhandled(ev)
		return
	}

	err := handler(ev.Data)

	if ev.ReplyTo.HasValue() {
		ev.ReplyTo.Value() <- err
	}
}

func (el *EventLoop) RegisterHandler(typ int, handler func(any) error) {
	el.handlers.Set(typ, handler)
}

func (el *EventLoop) SetUnhandledHandler(handler func(Event)) {
	el.handleUnhandled = handler
}

func (el *EventLoop) Post(ev Event) {
	el.events <- ev
}

func (el *EventLoop) PostAndWait(ev Event) any {
	reply := make(chan any, 1)
	ev.ReplyTo = MakeOptional(reply)
	el.Post(ev)
	return <-reply
}

func (el *EventLoop) Run() {
	ticker := time.NewTicker(time.Millisecond * 1) // Your "Resolution"
	defer ticker.Stop()

	for {
		select {
		case <-el.stop:
			return
		case ev := <-el.events:
			el.process(ev)
		case <-ticker.C:
			el.Post(Event{Type: 0})
		}
	}
}

func (el *EventLoop) PostAndGet(ev Event) Optional[any] {
	reply := make(chan any, 1)
	ev.ReplyTo = MakeOptional(reply)
	el.Post(ev)
	return MakeOptional(<-reply)
}

func (el *EventLoop) PostAndGetOk(ev Event) (any, bool) {
	reply := make(chan any, 1)
	ev.ReplyTo = MakeOptional(reply)
	el.Post(ev)
	val := <-reply
	return val, true
}

func (el *EventLoop) PostAndGetOrElse(ev Event, orElse func() any) any {
	reply := make(chan any, 1)
	ev.ReplyTo = MakeOptional(reply)
	el.Post(ev)
	select {
	case val := <-reply:
		return val
	default:
		return orElse()
	}
}
