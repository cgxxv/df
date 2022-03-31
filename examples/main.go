package main

import (
	"context"
	"fmt"
	"os"

	"github.com/cgxxv/df/dag"
)

func main() {
	dag.RegisterTask(&TaskDemo1{})
	dag.RegisterTask(&TaskDemo2{})
	dag.RegisterTask(&TaskDemo3{})
	dag.RegisterTask(&TaskDemo4{})
	dag.RegisterTask(&TaskDemo5{})

	var bus MyBus
	var ctx = context.TODO()
	var nms = dag.GetSchedulerNodeMeta(ctx, os.Getenv("SERVICE"))
	var service = "hello"
	var runTimes = 1

	if err := dag.InitScheduler(ctx, service); err != nil {
		panic(err)
	}

	dag.Schedule(ctx, service, &bus)

	if int(bus.Data) != runTimes*len(nms) {
		panic(fmt.Sprint("fuck ", bus.Data, runTimes*len(nms)))
	}

	fmt.Printf("%#v\n", bus)
}
