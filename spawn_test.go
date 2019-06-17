package gopactor

import (
	"testing"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/stretchr/testify/assert"
)

type TestActor struct{}

func (ta *TestActor) Receive(ctx actor.Context) {}

func TestSpawnFromInstance(t *testing.T) {
	a := assert.New(t)

	_, err := SpawnFromProducer(func () actor.Actor { return &TestActor{}} )
	a.Nil(err)

	// Cleanup
	PactReset()
}

func TestSpawnFromInstance_WithPrefix(t *testing.T) {
	a := assert.New(t)

	object, err := SpawnFromProducer(
		func () actor.Actor { return &TestActor{}},
		OptDefault.WithPrefix("test-actor"))

	a.Nil(err)
	a.Contains(object.String(), "test-actor")

	// Cleanup
	PactReset()
}

func TestSpawnFromProducer(t *testing.T) {
	a := assert.New(t)

	f := func() actor.Actor {
		return &TestActor{}
	}

	_, err := SpawnFromProducer(f)
	a.Nil(err)

	// Cleanup
	PactReset()
}

func TestSpawnFromFunc(t *testing.T) {
	a := assert.New(t)

	_, err := SpawnFromFunc(func(ctx actor.Context) {})
	a.Nil(err)

	// Cleanup
	PactReset()
}

func TestSpawnNullActor(t *testing.T) {
	a := assert.New(t)

	_, err := SpawnNullActor()
	a.Nil(err)

	// Cleanup
	PactReset()
}
