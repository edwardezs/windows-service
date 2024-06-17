package broadcaster

type Channel[T any] struct {
    bc broadcaster[*Channel[T]]
    ch chan T
    closed bool
}

func (ch *Channel[T]) Ch() <-chan T {
    return ch.ch
}

func (ch *Channel[T]) Get() T {
    return <-ch.ch
}

func (ch *Channel[T]) unregisterNoLock() {
    if ch.closed {
		return
	}
	ch.closed = true

    delete(ch.bc.getListeners(), ch)
	close(ch.ch)
}

func (ch *Channel[T]) Unregister() {
    if ch.closed {
		return
	}

    mutex := ch.bc.getMutex()
    mutex.Lock()
    defer mutex.Unlock()

    ch.unregisterNoLock()
}
