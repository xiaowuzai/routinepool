package main

import (
	"github.com/xiaowuzai/routinepool/demo1/pool1"
	"time"
)

func main(){
	p := pool1.New(5)

	for i := 0; i < 10; i++ {
		err := p.Schedule(func(){
				time.Sleep(time.Second * 3)
		})
		if err != nil {
			println("task: ", i, " err: ", err)
		}
	}

	p.Free()
}