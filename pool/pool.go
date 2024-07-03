package pool

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type Worker interface {
	Work(id int32)
}

type Pool struct {
	minRoutines       int32         // minimum number of routines that needs to be part of pool.
	pendingGoRoutines int32         // similar to above but no. of goroutines pending at a time.
	counter           int32         // number of goroutines that have ran and to give each routine a unique id.
	tasks             chan Worker   // channel in which task will be sent.
	shutdown          chan struct{} // to indicate that user has called to destroy this workpool.
	wg                sync.WaitGroup
}

func New(minRoutines int32) (*Pool, error) {
	if minRoutines < 0 {
		return nil, fmt.Errorf("not a valid number of goroutine, provided: %d", minRoutines)
	}

	pool := Pool{
		minRoutines: minRoutines,
		tasks:       make(chan Worker),
	}
	pool.Controller() // once the pool is initialized, start it. So that work can be submitted and processed.

	return &pool, nil
}

func (p *Pool) Register(work Worker) {
	atomic.AddInt32(&p.pendingGoRoutines, 1)
	p.tasks <- work
	// since it's an unbuffered channel, if we reach this point it means on pushing into channel was success and hence one task completed.
	atomic.AddInt32(&p.pendingGoRoutines, -1)
}

func (p *Pool) processTask(task Worker) {
	p.wg.Add(1)
	processId := p.counter
	fmt.Println("STARTING WORK", processId)
	defer p.wg.Done()
	task.Work(processId)
	fmt.Println("WORK DONE!!", processId)
}

func (p *Pool) Controller() {
	p.wg.Add(1) // since controller itself is a goroutine.
	fmt.Println("ADDING IN WORK WG")
	go func() {
		for {
			select {
			case task := <-p.tasks:
				fmt.Println("TASK RECEIVED")
				p.counter++ // increment counter one task is now received from channel.
				go p.processTask(task)
			case <-p.shutdown:
				fmt.Println("CLEANUP CODE HERE")
				p.wg.Done()
				close(p.tasks)
				close(p.shutdown)
			}
		}
	}()
}

func (p *Pool) Destroy() {
	// signal the shutdown phase
	fmt.Println("REMOVING IN WORK WG")
	fmt.Println("CLOSED CHANNELS")
	p.shutdown <- struct{}{}
}
