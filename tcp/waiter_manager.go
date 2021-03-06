package tcp

import (
	"sync"
	"time"
	log "github.com/sirupsen/logrus"
	"context"
)

type waiterManager struct {
	waiter map[int64]*waiter
	waiterLock *sync.RWMutex
	waiterGlobalTimeout int64 // 毫秒
	ctx context.Context
}

func newWaiterManager(ctx context.Context) *waiterManager {
	m := &waiterManager{
		waiter:              make(map[int64]*waiter),
		waiterLock:          new(sync.RWMutex),
		waiterGlobalTimeout: defaultWaiterTimeout,
		ctx:                 ctx,
	}
	go m.checkWaiterTimeout()
	return m
}

func (m *waiterManager) append(wai *waiter) {
	if m == nil {
		return
	}
	m.waiterLock.Lock()
	m.waiter[wai.msgId] = wai
	m.waiterLock.Unlock()
}

func (m *waiterManager) clear(msgId int64) {
	log.Infof("clear %v", msgId)
	if m == nil {
		return
	}
	m.waiterLock.Lock()
	wai, ok := m.waiter[msgId]
	if ok {
		delete(m.waiter, msgId)
		wai.reset()
	}
	m.waiterLock.Unlock()
}

func (m *waiterManager) clearTimeout() {
	if m == nil {
		return
	}
	current := int64(time.Now().UnixNano() / 1000000)
	m.waiterLock.Lock()
	for msgId, v := range m.waiter  {
		// check timeout
		if current - v.time >= m.waiterGlobalTimeout {
			log.Warnf("clearTimeout, msgid=[%v] is timeout, will delete", msgId)
			//close(v.data)
			delete(m.waiter, msgId)
			v.reset()
			v.StopWait()
			//tcp.wg.Done()
			// 这里为什么不能使用delWaiter的原因是
			// tcp.waiterLock已加锁，而delWaiter内部也加了锁
			// tcp.delWaiter(msgId)
		}
	}
	m.waiterLock.Unlock()
}

func (m *waiterManager) get(msgId int64) *waiter {
	if m == nil {
		return nil
	}
	m.waiterLock.RLock()
	wai, ok := m.waiter[msgId]
	m.waiterLock.RUnlock()
	if ok {
		return wai
	}
	log.Errorf("waiterManager get waiter not found, msgId=[%v]", msgId)
	return nil
}

func (m *waiterManager) clearAll() {
	m.waiterLock.Lock()
	for msgId, v := range m.waiter  {
		log.Infof("clearAll, %v stop wait", msgId)
		v.reset()
		v.StopWait()
		delete(m.waiter, msgId)
	}
	m.waiterLock.Unlock()
}

func (m *waiterManager) checkWaiterTimeout() {
	if m == nil {
		return
	}
	for {
		select {
		case <-m.ctx.Done():
			return
		default:
		}
		m.clearTimeout()
		time.Sleep(time.Second * 3)
	}
}

