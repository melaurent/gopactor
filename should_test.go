package gopactor

import (
	"testing"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/melaurent/gopactor/options"
	"github.com/stretchr/testify/assert"
)

func TestShouldReceive(t *testing.T) {
	a := assert.New(t)
	ctx := actor.EmptyRootContext

	receiver := SpawnFromFunc(func(ctx actor.Context) {}, OptDefault.WithPrefix("rcv"))

	// Wrong params
	a.Contains(ShouldReceive(nil), "not an actor PID")
	a.Contains(ShouldReceive(receiver), "One parameter with a message is required")

	// Failure: Timeout
	a.Contains(ShouldReceive(receiver, "Welcome"), "Timeout")

	// Failure: Message mismatch
	ctx.Send(receiver, "Hello, world!")
	a.Contains(ShouldReceive(receiver, "Welcome"), "do not match")

	// Success: Massage match
	ctx.Send(receiver,"Hello, world!")
	a.Empty(ShouldReceive(receiver, "Hello, world!"))

	// Success: Any message
	ctx.Send(receiver,"Hello, world!")
	a.Empty(ShouldReceive(receiver, nil))

	// Cleanup
	PactReset()
}

func TestShouldReceiveSomething(t *testing.T) {
	a := assert.New(t)
	ctx := actor.EmptyRootContext

	receiver := SpawnFromFunc(func(ctx actor.Context) {}, OptDefault.WithPrefix("rcv"))

	// Wrong params
	a.Contains(ShouldReceiveSomething(nil), "not an actor PID")

	// Failure: Timeout
	a.Contains(ShouldReceiveSomething(receiver), "Timeout")

	// Success: Any message
	ctx.Send(receiver, "Hello, world!")
	a.Empty(ShouldReceiveSomething(receiver))

	// Cleanup
	PactReset()
}

func TestShouldReceiveFrom(t *testing.T) {
	a := assert.New(t)
	ctx := actor.EmptyRootContext

	receiver := SpawnFromFunc(func(ctx actor.Context) {}, OptDefault.WithPrefix("rcv"))
	teller := SpawnFromFunc(func(ctx actor.Context) {}, OptDefault.WithPrefix("tel"))
	requestor := SpawnFromFunc(func(ctx actor.Context) {}, OptDefault.WithPrefix("req"))

	// Wrong params
	a.Contains(ShouldReceiveFrom(nil), "not an actor PID")
	a.Contains(ShouldReceiveFrom(receiver, teller), "Two parameters are required")
	a.Contains(ShouldReceiveFrom(receiver, nil, nil), "Sender should be an actor PID")

	// Failure: Timeout
	a.Contains(ShouldReceiveFrom(receiver, requestor, "from requestor"), "Timeout")

	// Failure: Sender unknown
	// When tell is used the receiver does not know who is the sender
	// NB: This protoactor behaviour might change in the future.
	// teller.Tell("ping")
	ctx.Send(receiver, "from teller")
	a.Contains(ShouldReceiveFrom(receiver, teller, "from teller"), "Sender is unknown")

	// Failure: Message mismatch
	ctx.RequestWithCustomSender(receiver, "from requestor", requestor)
	a.Contains(ShouldReceiveFrom(receiver, requestor, "from teller"), "Messages do not match")

	// Failure: Sender mismatch
	ctx.RequestWithCustomSender(receiver, "from requestor", requestor)
	a.Contains(ShouldReceiveFrom(receiver, teller, "from requestor"), "Sender does not match")

	// Success: everything matches
	ctx.RequestWithCustomSender(receiver, "from requestor", requestor)
	a.Empty(ShouldReceiveFrom(receiver, requestor, "from requestor"))

	// Success: any message
	ctx.RequestWithCustomSender(receiver, "from requestor", requestor)
	a.Empty(ShouldReceiveFrom(receiver, requestor, nil))

	// Cleanup
	PactReset()
}

func TestShouldReceiveN(t *testing.T) {
	a := assert.New(t)
	ctx := actor.EmptyRootContext

	receiver := SpawnFromFunc(func(ctx actor.Context) {}, OptDefault.WithPrefix("rcv"))

	// Wrong params
	a.Contains(ShouldReceiveN(nil), "not an actor PID")
	a.Contains(ShouldReceiveN(receiver), "the number of expected messages is required")
	a.Contains(ShouldReceiveN(receiver, 0), "should be a positive integer")
	a.Contains(ShouldReceiveN(receiver, "abc"), "should be a positive integer")

	// Failure: Nothing received at all
	a.Contains(ShouldReceiveN(receiver, 1), "got 0")

	// Failure: Not enough messages received

	ctx.Send(receiver, "Something")
	a.Contains(ShouldReceiveN(receiver, 2), "got 1")

	// Success: One message
	ctx.Send(receiver, "Something")
	a.Empty(ShouldReceiveN(receiver, 1))

	// Success: Many messages
	many := 30
	for i := 0; i < many; i++ {
		ctx.Send(receiver, i)
	}
	a.Empty(ShouldReceiveN(receiver, many))

	// Cleanup
	PactReset()
}

func TestShouldSend(t *testing.T) {
	a := assert.New(t)
	ctx := actor.EmptyRootContext

	receiver := ctx.SpawnPrefix(actor.PropsFromFunc(func(ctx actor.Context) {}), "rcv")
	sender := SpawnFromFunc(func(ctx actor.Context) {
		switch m := ctx.Message().(type) {
		case string:
			if m == "tell" {
				ctx.Send(receiver, "tell from sender")
			} else if m == "request" {
				ctx.Request(receiver, "request from sender")
			}
		}
	}, options.OptOutboundInterceptionOnly.WithPrefix("snd"))

	// Wrong params
	a.Contains(ShouldSend(nil), "not an actor PID")
	a.Contains(ShouldSend(sender), "One parameter with a message is required")

	// Failure: Timeout
	a.Contains(ShouldSend(sender, "from sender"), "Timeout")

	// Failure: Message mismatch
	ctx.Send(sender, "tell")
	a.Contains(ShouldSend(sender, "foobar"), "do not match")

	// Success: Tell: Massage match
	ctx.Send(sender, "tell")
	a.Empty(ShouldSend(sender, "tell from sender"))

	// Success: Tell: Any message
	ctx.Send(sender, "tell")
	a.Empty(ShouldSend(sender, nil))

	// Success: Request: Massage match
	ctx.Send(sender, "request")
	a.Empty(ShouldSend(sender, "request from sender"))

	// Success: Request: Any message
	ctx.Send(sender, "request")
	a.Empty(ShouldSend(sender, nil))

	// Cleanup
	PactReset()
}

func TestShouldSendTo(t *testing.T) {
	a := assert.New(t)
	ctx := actor.EmptyRootContext

	receiver := ctx.SpawnPrefix(actor.PropsFromFunc(func(ctx actor.Context) {}), "rcv")
	sender := SpawnFromFunc(func(ctx actor.Context) {
		switch m := ctx.Message().(type) {
		case string:
			if m == "tell" {
				ctx.Send(receiver, "tell from sender")
			} else if m == "request" {
				ctx.Request(receiver, "request from sender")
			}
		}
	}, options.OptOutboundInterceptionOnly.WithPrefix("snd"))

	// Wrong params
	a.Contains(ShouldSendTo(nil), "not an actor PID")
	a.Contains(ShouldSendTo(sender), "Two parameters are required")
	a.Contains(ShouldSendTo(sender, nil), "Two parameters are required")
	a.Contains(ShouldSendTo(sender, nil, nil), "Receiver should be an actor PID")

	// Failure: Timeout
	a.Contains(ShouldSendTo(sender, receiver, "from sender"), "Timeout")

	// Failure: Message mismatch
	ctx.Send(sender, "tell")
	a.Contains(ShouldSendTo(sender, receiver, "foobar"), "do not match")

	// Failure: Receiver mismatch
	ctx.Send(sender, "tell")
	a.Contains(ShouldSendTo(sender, sender, "tell from sender"), "Receiver does not match")

	// Success: Tell: Massage match
	ctx.Send(sender, "tell")
	a.Empty(ShouldSendTo(sender, receiver, "tell from sender"))

	// Success: Tell: Any message
	ctx.Send(sender, "tell")
	a.Empty(ShouldSendTo(sender, receiver, nil))

	// Success: Request: Massage match
	ctx.Send(sender, "request")
	a.Empty(ShouldSendTo(sender, receiver, "request from sender"))

	// Success: Request: Any message
	ctx.Send(sender, "request")
	a.Empty(ShouldSendTo(sender, receiver, nil))

	// Cleanup
	PactReset()
}

func TestShouldSendSomething(t *testing.T) {
	a := assert.New(t)
	ctx := actor.EmptyRootContext

	receiver := ctx.SpawnPrefix(actor.PropsFromFunc(func(ctx actor.Context) {}), "rcv")
	sender := SpawnFromFunc(func(ctx actor.Context) {
		switch m := ctx.Message().(type) {
		case string:
			if m == "tell" {
				ctx.Send(receiver, "tell from sender")
			} else if m == "request" {
				ctx.Request(receiver, "request from sender")
			}
		}
	}, options.OptOutboundInterceptionOnly.WithPrefix("snd"))

	// Wrong params
	a.Contains(ShouldSendSomething(nil), "not an actor PID")

	// Failure: Timeout
	a.Contains(ShouldSendSomething(sender), "Timeout")

	// Success: Tell: Any message
	ctx.Send(sender, "tell")
	a.Empty(ShouldSendSomething(sender))

	// Success: Request: Any message
	ctx.Send(sender, "request")
	a.Empty(ShouldSendSomething(sender))

	// Cleanup
	PactReset()
}

func TestShouldSendN(t *testing.T) {
	a := assert.New(t)
	ctx := actor.EmptyRootContext

	receiver := ctx.SpawnPrefix(actor.PropsFromFunc(func(ctx actor.Context) {}), "rcv")
	sender := SpawnFromFunc(func(ctx actor.Context) {
		switch m := ctx.Message().(type) {
		case string:
			if m == "tell" {
				ctx.Send(receiver, "tell from sender")
			} else if m == "request" {
				ctx.Request(receiver, "request from sender")
			}
		}
	}, options.OptOutboundInterceptionOnly.WithPrefix("snd"))

	// Wrong params
	a.Contains(ShouldSendN(nil), "not an actor PID")
	a.Contains(ShouldSendN(sender), "number of expected messages is required")
	a.Contains(ShouldSendN(sender, 0), "should be a positive integer")
	a.Contains(ShouldSendN(sender, "abc"), "should be a positive integer")

	// Failure: Not sending at all
	a.Contains(ShouldSendN(sender, 1), "got 0")

	// Failure: Tell: Not sending enough
	ctx.Send(sender, "tell")
	a.Contains(ShouldSendN(sender, 2), "got 1")

	// Success: Tell: one message
	ctx.Send(sender, "tell")
	a.Empty(ShouldSendN(sender, 1))

	// Success: Request: one message
	ctx.Send(sender, "request")
	a.Empty(ShouldSendN(sender, 1))

	// Success: Many messages
	many := 30
	for i := 0; i < many; i++ {
		ctx.Send(sender, "tell")
	}
	a.Empty(ShouldSendN(sender, many))

	// Cleanup
	PactReset()
}

func TestShouldNotSendOrReceive(t *testing.T) {
	a := assert.New(t)
	ctx := actor.EmptyRootContext

	receiver := SpawnFromFunc(func(ctx actor.Context) {}, OptDefault.WithPrefix("rcv"))
	sender := SpawnFromFunc(func(ctx actor.Context) {
		switch m := ctx.Message().(type) {
		case string:
			if m == "tell" {
				ctx.Send(receiver, "tell from sender")
			} else if m == "request" {
				ctx.Request(receiver, "request from sender")
			}
		}
	}, options.OptOutboundInterceptionOnly.WithPrefix("snd"))

	// Wrong params
	a.Contains(ShouldNotSendOrReceive(nil), "not an actor PID")

	// Success: neither send, nor receive
	a.Empty(ShouldNotSendOrReceive(sender))

	// Failure: receive
	ctx.Send(receiver, "foobar")
	a.Contains(ShouldNotSendOrReceive(receiver), "Got inbound message")

	// Failure: tell
	ctx.Send(sender, "tell")
	a.Contains(ShouldNotSendOrReceive(sender), "Got outbound message")

	// Failure: request
	ctx.Send(sender, "request")
	a.Contains(ShouldNotSendOrReceive(sender), "Got outbound message")

	// Cleanup
	PactReset()
}

func TestShouldStart(t *testing.T) {
	a := assert.New(t)

	receiver := SpawnFromFunc(func(ctx actor.Context) {}, OptNoInterception.WithSystemInterception().WithPrefix("rcv"))

	// Wrong params
	a.Contains(ShouldStart(nil), "not an actor PID")

	// Success
	a.Empty(ShouldStart(receiver))

	// Cleanup
	PactReset()
}

func TestShouldStop(t *testing.T) {
	a := assert.New(t)

	receiver := SpawnFromFunc(func(ctx actor.Context) {}, OptNoInterception.WithSystemInterception().WithPrefix("rcv"))

	// Wrong params
	a.Contains(ShouldStop(nil), "not an actor PID")

	// Failure: Timeout
	a.Contains(ShouldStop(receiver), "Timeout")

	// Success
	receiver.Stop()
	a.Empty(ShouldStop(receiver))

	// Cleanup
	PactReset()
}

func TestShouldBeRestarting(t *testing.T) {
	a := assert.New(t)
	ctx := actor.EmptyRootContext

	receiver := SpawnFromFunc(func(ctx actor.Context) {
		switch m := ctx.Message().(type) {
		case string:
			if m == "panic" {
				panic("I am panicing!")
			}
		}
	}, OptNoInterception.WithSystemInterception().WithPrefix("rcv"))

	// Wrong params
	a.Contains(ShouldBeRestarting(nil), "not an actor PID")

	// Failure: Timeout
	a.Contains(ShouldBeRestarting(receiver), "Timeout")

	// Success
	ctx.Send(receiver, "panic")
	a.Empty(ShouldBeRestarting(receiver))
	a.Empty(ShouldStart(receiver))

	// Cleanup
	PactReset()
}

func TestShouldObserveTermination(t *testing.T) {
	a := assert.New(t)

	childProps := actor.PropsFromFunc(func(ctx actor.Context) {})

	getParentChildSet := func() (*actor.PID, *actor.PID) {
		var child *actor.PID
		wait := make(chan bool)
		parent := SpawnFromFunc(func(ctx actor.Context) {
			switch ctx.Message().(type) {
			case *actor.Started:
				child = ctx.SpawnPrefix(childProps, "child")
				wait <- true
			}
		}, options.OptNoInterception.WithSystemInterception().WithPrefix("parent"))

		<-wait

		return parent, child
	}

	someActor := SpawnNullActor()

	// Wrong params
	parent1, child1 := getParentChildSet()
	a.Contains(ShouldObserveTermination(nil), "not an actor PID")
	a.Contains(ShouldObserveTermination(parent1, nil), "should be an actor PID")

	// Failure: Timeout
	a.Contains(ShouldObserveTermination(parent1, child1), "Timeout")

	// Failure: pid mismatch
	child1.Stop()
	a.Contains(ShouldObserveTermination(parent1, someActor), "Timeout")

	// Success: Termination of the child
	parent2, child2 := getParentChildSet()
	child2.Stop()
	a.Empty(ShouldObserveTermination(parent2, child2))

	// Success: Any termination
	parent3, child3 := getParentChildSet()
	child3.Stop()
	a.Empty(ShouldObserveTermination(parent3))

	// Cleanup
	PactReset()
}

func TestShouldSpawn(t *testing.T) {
	a := assert.New(t)
	ctx := actor.EmptyRootContext

	receiver := SpawnFromFunc(func(ctx actor.Context) {}, OptNoInterception.WithPrefix("rcv"))
	childProps := actor.PropsFromFunc(func(ctx actor.Context) {})
	parent := SpawnFromFunc(func(ctx actor.Context) {
		switch m := ctx.Message().(type) {
		case string:
			if m == "spawn" && ctx.Sender() != nil {
				child := ctx.SpawnPrefix(childProps, "my-dear-dummy")
				ctx.Respond(child)
			}
		}
	}, options.OptNoInterception.WithSpawnInterception().WithPrefix("parent"))

	// Wrong params
	a.Contains(ShouldSpawn(nil), "not an actor PID")
	a.Contains(ShouldSpawn(parent, 123), "should be a string")

	// Failure: Timeout
	a.Contains(ShouldSpawn(parent), "Timeout")

	// Failure: child PID mismatch
	ctx.RequestWithCustomSender(parent, "spawn", receiver)
	a.Contains(ShouldSpawn(parent, "foobar"), "does not match")

	// Success: child PID match
	ctx.RequestWithCustomSender(parent, "spawn", receiver)
	a.Empty(ShouldSpawn(parent, "my-dear-dummy"))

	// Success: any child PID is ok
	ctx.RequestWithCustomSender(parent, "spawn", receiver)
	a.Empty(ShouldSpawn(parent))

	// Cleanup
	PactReset()
}
