package gopactor

import (
	"testing"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/stretchr/testify/assert"
)

type TestActor struct{}

func (ta *TestActor) Receive(ctx actor.Context) {}

func TestSpawnFromInstance_WithPrefix(t *testing.T) {
	a := assert.New(t)

	object := SpawnFromProducer(
		func () actor.Actor { return &TestActor{}},
		OptDefault.WithPrefix("test-actor"))

	a.Contains(object.String(), "test-actor")

	// Cleanup
	PactReset()
}