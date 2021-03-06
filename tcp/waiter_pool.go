package tcp

import (
	"sync"
	"time"
)

type waiterPool struct {
	maxSize int64
	pool []*waiter
	lock *sync.Mutex
}

func newWaiterPool(maxSize int64) *waiterPool {
	if maxSize < 8 {
		maxSize = 8
	}
	p := &waiterPool{
		maxSize: maxSize,
		pool: make([]*waiter, 64),
		lock: new(sync.Mutex),
	}
	for i := 0; i < 64; i++ {
		p.pool[i] = newWaiter(0, nil)
	}
	return p
}

func (p *waiterPool) get(msgId int64, oncomplete func(i int64)) (*waiter, error) {
	p.lock.Lock()
	defer p.lock.Unlock()
	for _, w := range p.pool {
		if w.msgId <= 0 {
			w.msgId = msgId
			w.onComplete = oncomplete
			w.time = int64(time.Now().UnixNano() / 1000000)
			return w, nil
		}
	}
	if int64(len(p.pool)) < p.maxSize - 1 {
		w := newWaiter(msgId, oncomplete)
		p.pool = append(p.pool, w)
		return w, nil
	}
	return nil, ErrPoolMaxSize
}
