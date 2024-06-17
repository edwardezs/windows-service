package broadcaster

type Listener[T any] struct {
	bc broadcaster[*Listener[T]]
	ch chan Payload[T]
	closed bool
}

func (l *Listener[T]) Ch() <-chan Payload[T] {
	return l.ch
}

func (l *Listener[T]) Get() T {
	payload := <-l.Ch()
	return payload.Get()
}

func (l *Listener[T]) unregisterNoLock() {
	if l.closed {
		return
	}
	l.closed = true

	delete(l.bc.getListeners(), l)
	close(l.ch)
}

func (l *Listener[T]) Unregister() {
	if l.closed {
		return
	}

	mutex := l.bc.getMutex()
    mutex.Lock()
    defer mutex.Unlock()
	
	l.unregisterNoLock()
}
