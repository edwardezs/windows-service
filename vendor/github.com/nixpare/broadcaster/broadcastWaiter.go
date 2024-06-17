package broadcaster

import "sync"

type BroadcastWaiter[T any] struct {
	listeners map[*Listener[T]]struct{}
	lMutex    *sync.RWMutex
	closed    bool
}

func NewBroadcastWaiter[T any]() *BroadcastWaiter[T] {
	return &BroadcastWaiter[T]{
		listeners: make(map[*Listener[T]]struct{}),
		lMutex:    new(sync.RWMutex),
	}
}

func (bc *BroadcastWaiter[T]) getListeners() map[*Listener[T]]struct{} {
	return bc.listeners
}

func (bc *BroadcastWaiter[T]) getMutex() *sync.RWMutex {
	return bc.lMutex
}

func (bc *BroadcastWaiter[T]) Register(bufSize int) *Listener[T] {
	bc.lMutex.Lock()
	defer bc.lMutex.Unlock()

	l := &Listener[T]{
		bc: bc,
		ch: make(chan Payload[T], bufSize),
	}

	if bc.closed {
		l.closed = true
		close(l.ch)
	} else {
		bc.listeners[l] = struct{}{}
	}

	return l
}

func (bc *BroadcastWaiter[T]) Send(payload T) Waiter[T] {
	bc.lMutex.RLock()
	defer bc.lMutex.RUnlock()

	w := &broadcastPayloadWaiter[T]{
		listeners: make([]*Listener[T], 0, len(bc.listeners)),
		payload:   payload,
	}
	w.wg.Add(len(bc.listeners))

	for l := range bc.listeners {
		w.listeners = append(w.listeners, l)
		l.ch <- w
	}

	return w
}

func (bc *BroadcastWaiter[T]) Close() {
	bc.lMutex.Lock()
	defer bc.lMutex.Unlock()

	if bc.closed {
		return
	}
	bc.closed = true

	for l := range bc.listeners {
		l.unregisterNoLock()
	}
}

type BufBroadcastWaiter[T any] struct {
	listeners map[*Listener[T]]struct{}
	lMutex    *sync.RWMutex
	closed    bool
	data      []T
}

func NewBufBroadcastWaiter[T any]() *BufBroadcastWaiter[T] {
	return &BufBroadcastWaiter[T]{
		listeners: make(map[*Listener[T]]struct{}),
		lMutex:    new(sync.RWMutex),
	}
}

func (bc *BufBroadcastWaiter[T]) getListeners() map[*Listener[T]]struct{} {
	return bc.listeners
}

func (bc *BufBroadcastWaiter[T]) getMutex() *sync.RWMutex {
	return bc.lMutex
}

func (bc *BufBroadcastWaiter[T]) registerNoLock(bufSize int) *Listener[T] {
	l := &Listener[T]{
		bc: bc,
		ch: make(chan Payload[T], bufSize),
	}

	if bc.closed {
		l.closed = true
		close(l.ch)
	} else {
		bc.listeners[l] = struct{}{}
	}

	bc.listeners[l] = struct{}{}
	return l
}

func (bc *BufBroadcastWaiter[T]) Register(bufSize int) *Listener[T] {
	bc.lMutex.Lock()
	defer bc.lMutex.Unlock()

	return bc.registerNoLock(bufSize)
}

func (bc *BufBroadcastWaiter[T]) Send(payload T) Waiter[T] {
	bc.lMutex.RLock()
	defer bc.lMutex.RUnlock()

	bc.data = append(bc.data, payload)

	w := &broadcastPayloadWaiter[T]{
		listeners: make([]*Listener[T], 0, len(bc.listeners)),
		payload:   payload,
	}
	w.wg.Add(len(bc.listeners))

	for l := range bc.listeners {
		w.listeners = append(w.listeners, l)
		l.ch <- w
	}

	return w
}

func (bc *BufBroadcastWaiter[T]) Data() []T {
	return bc.data
}

func (bc *BufBroadcastWaiter[T]) Connect(bufSize int) ([]T, *Listener[T]) {
	bc.lMutex.Lock()
	defer bc.lMutex.Unlock()

	return bc.Data(), bc.registerNoLock(bufSize)
}

func (bc *BufBroadcastWaiter[T]) Close() {
	bc.lMutex.Lock()
	defer bc.lMutex.Unlock()

	if bc.closed {
		return
	}
	bc.closed = true

	for l := range bc.listeners {
		l.unregisterNoLock()
	}
}
