package broadcaster

type Receiver[T any] chan Payload[T]

func NewReceiver[T any](size int) *Receiver[T] {
    r := new(Receiver[T])
    *r = make(chan Payload[T])
    return r
}

func (r *Receiver[T]) Send(payload T) Waiter[T] {
    pw := &receivePayloadWaiter[T]{
        data: payload,
        ch:   make(chan struct{}),
    }

    *r <- pw
    return pw
}

func (r *Receiver[T]) Ch() <-chan Payload[T] {
    return *r
}

func (r *Receiver[T]) Close() {
    close(*r)
}