package broadcaster

import (
	"sync"
)

type broadcaster[T comparable] interface {
    getListeners() map[T]struct{}
    getMutex()     *sync.RWMutex
}

type Broadcaster[T any] struct {
    listeners map[*Channel[T]]struct{}
    lMutex    *sync.RWMutex
    closed    bool
}

func NewBroadcaster[T any]() *Broadcaster[T] {
    return &Broadcaster[T]{
        listeners: make(map[*Channel[T]]struct{}),
        lMutex:    new(sync.RWMutex),
    }
}

func (bc *Broadcaster[T]) getListeners() map[*Channel[T]]struct{} {
    return bc.listeners
}

func (bc *Broadcaster[T]) getMutex() *sync.RWMutex {
    return bc.lMutex
}

func (bc *Broadcaster[T]) Register(bufSize int) *Channel[T] {
    bc.lMutex.Lock()
    defer bc.lMutex.Unlock()

    ch := &Channel[T]{ bc: bc, ch: make(chan T, bufSize) }
    
    if bc.closed {
        ch.closed = true
        close(ch.ch)
    } else {
        bc.listeners[ch] = struct{}{}
    }
   
    return ch
}

func (bc *Broadcaster[T]) Send(payload T) {
    bc.lMutex.RLock()
    defer bc.lMutex.RUnlock()

    for ch := range bc.listeners {
        ch.ch <- payload
    }
}

func (bc *Broadcaster[T]) Close() {
    bc.lMutex.Lock()
    defer bc.lMutex.Unlock()

    if bc.closed {
        return
    }
    bc.closed = true

    for ch := range bc.listeners {
        ch.unregisterNoLock()
    }
}

type BufBroadcaster[T any] struct {
    listeners map[*Channel[T]]struct{}
    lMutex    *sync.RWMutex
    closed    bool
    data      []T
}

func NewBufBroadcaster[T any]() *BufBroadcaster[T] {
    return &BufBroadcaster[T]{
        listeners: make(map[*Channel[T]]struct{}),
        lMutex:    new(sync.RWMutex),
    }
}

func (bc *BufBroadcaster[T]) getListeners() map[*Channel[T]]struct{} {
    return bc.listeners
}

func (bc *BufBroadcaster[T]) getMutex() *sync.RWMutex {
    return bc.lMutex
}

func (bc *BufBroadcaster[T]) registerNoLock(bufSize int) *Channel[T] {
    ch := &Channel[T]{ bc: bc, ch: make(chan T, bufSize) }

    if bc.closed {
        ch.closed = true
        close(ch.ch)
    } else {
        bc.listeners[ch] = struct{}{}
    }
    
    return ch
}

func (bc *BufBroadcaster[T]) Register(bufSize int) *Channel[T] {
    bc.lMutex.Lock()
    defer bc.lMutex.Unlock()
    
    return bc.registerNoLock(bufSize)
}

func (bc *BufBroadcaster[T]) Send(payload T) {
    bc.lMutex.RLock()
    defer bc.lMutex.RUnlock()

    bc.data = append(bc.data, payload)
    for ch := range bc.listeners {
        ch.ch <- payload
    }
}

func (bc *BufBroadcaster[T]) Data() []T {
    return bc.data
}

func (bc *BufBroadcaster[T]) Connect(bufSize int) ([]T, *Channel[T]) {
    bc.lMutex.Lock()
    defer bc.lMutex.Unlock()

    return bc.Data(), bc.registerNoLock(bufSize)
}

func (bc *BufBroadcaster[T]) Close() {
    bc.lMutex.Lock()
    defer bc.lMutex.Unlock()

    if bc.closed {
        return
    }
    bc.closed = true

    for ch := range bc.listeners {
        ch.unregisterNoLock()
    }
}
