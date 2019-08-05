package ConcurrencyCron

import "errors"

/**
 *@author  wxn
 *@project ConcurrencyCron
 *@package ConcurrencyCron
 *@date    19-8-2 上午10:02
 */
var MAX_POOL_CAPITION uint32 = 10000

//Control the number of concurrent ticket pools
type TicketsPool interface {
	//take a ticket
	Take()
	//return a ticket
	Return()
	//get the ticket pools' status
	Active() bool
	//total tickets
	Total() uint32
	//remain tickets
	Remain() uint32
	//close pool
	Close()
}

type tickets struct {
	total  uint32        //total tickets
	ticket chan struct{} //ticket
	active bool          //Whether the ticket pool has been activated
}

//Set the maximum number of concurrent
func SetMaxConcurrent(max uint32) (err error) {
	if max == 0 {
		return errors.New("max concurrent must >0")
	}
	MAX_POOL_CAPITION = max
	return
}

//create a tickets pool
func NewTicketsPool(total uint32) (TicketsPool, error) {
	tp := tickets{}
	if err := tp.init(total); err != nil {
		return nil, err
	}
	return &tp, nil
}

func (tp *tickets) init(total uint32) (err error) {
	if tp.active {
		return errors.New("tickets pool is already active")
	}
	if total == 0 {
		return errors.New("tickets num must >0")
	}
	if total > MAX_POOL_CAPITION {
		total = MAX_POOL_CAPITION
	}
	ch := make(chan struct{}, total)
	n := int(total)
	for i := 0; i < n; i++ {
		ch <- struct{}{}
	}
	tp.ticket = ch
	tp.total = total
	tp.active = true
	return err
}

func (tp *tickets) Take() {
	<-tp.ticket
}

func (tp *tickets) Return() {
	tp.ticket <- struct{}{}
}

func (tp *tickets) Active() bool {
	return tp.active
}

func (tp *tickets) Total() uint32 {
	return tp.total
}

func (tp *tickets) Remain() uint32 {
	return uint32(len(tp.ticket))
}

func (tp *tickets) Close() {
	tp.active = false
	close(tp.ticket)
}
