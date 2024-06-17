package broadcaster

import "sync"

type Payload[T any] interface {
	Done()
	Data() T
	Get() T
}

type Waiter[T any] interface {
	Wait()
}

type broadcastPayloadWaiter[T any] struct {
	listeners []*Listener[T]
	payload   T
	wg        sync.WaitGroup
}

func (pw *broadcastPayloadWaiter[T]) Wait() {
	pw.wg.Wait()
}

func (pw *broadcastPayloadWaiter[T]) Done() {
	pw.wg.Done()
}

func (pw *broadcastPayloadWaiter[T]) Data() T {
	return pw.payload
}

func (pw *broadcastPayloadWaiter[T]) Get() T {
	defer pw.Done()
	return pw.Data()
}

type receivePayloadWaiter[T any] struct {
    data T
    ch chan struct{}
}

func (p *receivePayloadWaiter[T]) Data() T {
    return p.data
}

func (p *receivePayloadWaiter[T]) Done() {
    close(p.ch)
}

func (p *receivePayloadWaiter[T]) Get() T {
    defer p.Done()
    return p.data
}

func (p *receivePayloadWaiter[T]) Wait() {
    for range p.ch {}
}
