package pool1

import (
	"errors"
	"fmt"
	"sync"
)

type Task func()

type Pool struct {
	capacity int  // pool cap
	active chan struct{}  //   pool active chan

	tasks chan Task  // task chan

	wg sync.WaitGroup  // 在 pool 销毁时等待所有 goroutine 退出
	quit chan struct{}  // 通知 goroutine 退出的信号  channel
}

var (
	defaultCap int = 10
	maxCap int = 20
)

// 指定 active chan 的容量，并通过 goroutine 来执行 pool
func New(capacity int) *Pool {
	if capacity <= 0  {
		capacity = defaultCap
	}else if capacity > maxCap {
		capacity = maxCap
	}

	p := &Pool{
		capacity: capacity,
		active: make(chan struct{},capacity),
		tasks: make(chan Task),
		quit: make(chan struct{}),
	}

	fmt.Println("workpool start")
	go p.run()

	return p
}

//
func (p *Pool)run() {
	// 索引
	idx := 0
	for  {
		select {
		case <- p.quit:
			return
		case p.active <- struct{}{}:  // p.active 没有满
			idx++
			p.newWorker(idx)
		}
	}
}

func (p *Pool)newWorker(idx int) {
	p.wg.Add(1)

	go func(){
		defer func(){
			if err := recover(); err != nil {
				fmt.Printf("worker[%03d]: recover panic[%s] and exist\n", idx, err)
				<-p.active
			}
			p.wg.Done()
		}()

		fmt.Printf("worker[%03d]: start\n", idx)
		for  {
			select {
			case <-p.quit:
				fmt.Printf("worker[%03d]: exit\n", idx)
				<-p.active
				return
			case t := <-p.tasks:
				fmt.Printf("worker[%03d]: receive a task\n", idx)
				t()
			}
		}
	}()
}

var ErrWorkerPoolFreed =  errors.New("workerpool freed")

//Schedule 向 task channel 里面添加 task
func (p *Pool) Schedule(t Task) error {
	select {
	case <-p.quit:
		return ErrWorkerPoolFreed
	case p.tasks <- t:  //
		return nil
	}
}


func (p *Pool)Free() {
	close(p.quit) // 关闭 quit channel 。 <-p.quit 不会发生阻塞，起到通知的作用。
	p.wg.Wait()
	fmt.Println("workpool freed")
}



















