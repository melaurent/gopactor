package gopactor

import (
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
)

type Worker struct{}

func (w *Worker) Receive(ctx actor.Context) {}

func ExampleSpawnFromInstance() {
	ctx := actor.EmptyRootContext

	// Given that the Worker actor is defined elsewhere
	worker := SpawnFromProducer(func () actor.Actor { return &Worker{} })

	ctx.Send(worker, "Hello, world!")
}

func ExampleSpawnFromProducer() {
	ctx := actor.EmptyRootContext

	producer := func() actor.Actor {
		return &Worker{}
	}

	worker := SpawnFromProducer(producer)
	ctx.Send(worker, "Hello, world!")
}

func ExampleSpawnFromFunc() {
	ctx := actor.EmptyRootContext

	f := func(ctx actor.Context) {
		if msg, ok := ctx.Message().(string); ok {
			fmt.Printf("Got a message: %s\n", msg)
		}
	}

	worker := SpawnFromFunc(f)

	ctx.Send(worker,"Hello, world!")
	ShouldReceiveSomething(worker)
	// Output: Got a message: Hello, world!
}

func ExampleSpawnNullActor() {
	ctx := actor.EmptyRootContext

	worker := SpawnFromProducer(func () actor.Actor { return &Worker{} })
	requestor := SpawnNullActor()

	ctx.RequestWithCustomSender(worker, "ping", requestor)
}
