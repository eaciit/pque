package pque

import (
	"sync"
	"time"
)

type Que struct {
	WorkerCount, JobCount, PreparedJob int
	ProcessedJob, CompletedJob         int
	AllKeyHasBeenSent                  bool
	Fn                                 func(interface{}) interface{}
	FnDone                             func(interface{})
	FnWaiting                          func(*Que)
	IsRunning, Completed               bool

	keys           chan interface{}
	results        chan interface{}
	keySentChannel chan bool

	wg *sync.WaitGroup
}

func NewQue() *Que {
	q := new(Que)
	return q
}

func (q *Que) WaitForKeys() {
	q.initChannel()
	q.runProcess()
	go func() {
		for !q.AllKeyHasBeenSent {
			select {
			case <-q.keySentChannel:
				q.AllKeyHasBeenSent = true
			}
		}
	}()
}

func (q *Que) SendKey(k interface{}) {
	q.initChannel()
	if q.keys == nil {
		q.keys = make(chan interface{})
	}
	q.keys <- k
}

func (q *Que) KeySendDone() {
	q.initChannel()
	q.keySentChannel <- true
}

func (q *Que) WaitForCompletion() {
	q.initChannel()
	for !q.AllKeyHasBeenSent {
		time.Sleep(1 * time.Microsecond)
	}
	q.wg.Wait()
	q.IsRunning = false
	q.Completed = true
}

func (q *Que) initChannel() {
	if q.keys == nil {
		q.keys = make(chan interface{})
	}

	if q.results == nil {
		q.results = make(chan interface{})
	}

	if q.keySentChannel == nil {
		q.keySentChannel = make(chan bool)
	}

	if q.wg == nil {
		q.wg = new(sync.WaitGroup)
	}
}

func (q *Que) runProcess() {
	q.IsRunning = true
	if q.FnWaiting != nil {
		go q.FnWaiting(q)
	}
	for widx := 0; widx < q.WorkerCount; widx++ {
		go func(que *Que) {
			for !que.AllKeyHasBeenSent {
				for k := range que.keys {
					que.wg.Add(1)
					que.ProcessedJob++
					if que.Fn != nil {
						que.results <- que.Fn(k)
					}
					q.wg.Done()
				}
			}
		}(q)
	}

	go q.receiveResult()
}

func (q *Que) receiveResult() {
	q.initChannel()
	for q.IsRunning == true {
		select {
		case r := <-q.results:
			q.CompletedJob++
			if q.FnDone != nil {
				q.FnDone(r)
			}
		}
	}
}
