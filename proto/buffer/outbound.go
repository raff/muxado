package buffer

import (
	"sync"
)

type Outbound struct {
	val int
	err error
	*sync.Cond
}

func NewOutbound(size int) *Outbound {
	return &Outbound{val: size, Cond: sync.NewCond(new(sync.Mutex))}
}

func (b *Outbound) Increment(inc int) {
	b.L.Lock()
	defer b.L.Unlock()

	b.val += inc
	b.Broadcast()
}

func (b *Outbound) SetError(err error) {
	b.L.Lock()
	defer b.L.Unlock()

	b.err = err
	b.Broadcast()
}

func (b *Outbound) Decrement(dec int) (ret int, err error) {
	if dec == 0 {
		return
	}

	b.L.Lock()
	defer b.L.Unlock()

	for {
		if b.err != nil {
			err = b.err
			break
		}

		if b.val > 0 {
			if dec > b.val {
				ret = b.val
				b.val = 0
				break
			} else {
				b.val -= dec
				ret = dec
				break
			}
		} else {
			b.Wait()
		}
	}
	return
}
