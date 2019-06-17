package gopactor

import (
	"fmt"
	"log"

	"github.com/AsynkronIT/protoactor-go/actor"
)

type Worker struct{}

func (w *Worker) Receive(ctx actor.Context) {}

func ExampleSpawnFromInstance() {
	ctx := actor.EmptyRootContext

	// Given that the Worker actor is defined elsewhere
	worker, err := SpawnFromProducer(func () actor.Actor { return &Worker{} })
	if err != nil {
		log.Print("Failed to spawn a worker")
		return
	}

	ctx.Send(worker, "Hello, world!")
}

func ExampleSpawnFromProducer() {
	ctx := actor.EmptyRootContext

	producer := func() actor.Actor {
		return &Worker{}
	}

	worker, _ := SpawnFromProducer(producer)
	ctx.Send(worker, "Hello, world!")
}

func ExampleSpawnFromFunc() {
	ctx := actor.EmptyRootContext

	f := func(ctx actor.Context) {
		if msg, ok := ctx.Message().(string); ok {
			fmt.Printf("Got a message: %s\n", msg)
		}
	}

	worker, _ := SpawnFromFunc(f)

	ctx.Send(worker,"Hello, world!")
	ShouldReceiveSomething(worker)
	// Output: Got a message: Hello, world!
}

func ExampleSpawnNullActor() {
	ctx := actor.EmptyRootContext

	worker, _ := SpawnFromProducer(func () actor.Actor { return &Worker{} })
	requestor, _ := SpawnNullActor()

	ctx.RequestWithCustomSender(worker, "ping", requestor)
}
