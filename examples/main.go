package main

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"sync"

	"github.com/cgxxv/df/dag"
)

func main() {
	bollet := make([]byte, 10<<30)

	dag.RegisterTask(&TaskDemo1{})
	dag.RegisterTask(&TaskDemo2{})
	dag.RegisterTask(&TaskDemo3{})
	dag.RegisterTask(&TaskDemo4{})
	dag.RegisterTask(&TaskDemo5{})

	var bus MyBus
	var ctx = context.TODO()
	var nms = dag.GetSchedulerNodeinfo(ctx, os.Getenv("SERVICE"))
	var service = "hello"
	var runTimes = 1000

	if err := dag.InitScheduler(ctx, service); err != nil {
		panic(err)
	}

	var wg sync.WaitGroup
	for i := 0; i < runTimes; i++ {
		wg.Add(1)
		go func() {
			if err := dag.Schedule(ctx, service, &bus); err != nil {
				println(err.Error())
			}
			wg.Done()
		}()
	}
	wg.Wait()

	if int(bus.Data) != runTimes*len(nms) {
		println("NOT OK")
	}
	fmt.Printf("OK: %#v\n", bus)

	runtime.KeepAlive(bollet)
}
