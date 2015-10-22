package pque

import (
	"sync"
)

type Que struct {
	WorkerCount, JobCount, PreparedJob int
	ProcessedJob, CompletedJob         int
	SimultanProcess, AllKeyHasBeenSent bool
	Fn                                 func(in interface{}) interface{}
	FnDone                             func(in interface{})
	IsRunning                          bool

	keys           chan interface{}
	results        chan interface{}
	keySentChannel chan bool

	wg *sync.WaitGroup
}

func NewQue() *Que {
	q := new(Que)
	q.SimultanProcess = true
	q.wg = new(sync.WaitGroup)
	return q
}

func (q *Que) WaitForKeys() {
	q.initChannel()
	if q.SimultanProcess == true {
		q.runProcess()
	}
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
	if q.SimultanProcess == false {
		q.runProcess()
	}
}

func (q *Que) WaitForCompletion() {
	q.initChannel()
	for !q.AllKeyHasBeenSent {
		select {
		case <-q.keySentChannel:
			q.AllKeyHasBeenSent = true
		}
	}
	q.wg.Wait()
	q.IsRunning = false
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
}

func (q *Que) runProcess() {
	q.IsRunning = true
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
