package pact

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/meamidos/pact/catcher"
)

func (p *Pact) spawn(props *actor.Props, options ...catcher.Options) (*actor.PID, error) {
	catcher := catcher.New()
	catcher.LoggingOn = p.LoggingOn

	pid, err := catcher.Spawn(props, options...)
	if err != nil {
		return nil, err
	}

	p.CatchersByPID[pid.String()] = catcher

	return pid, nil
}

func (p *Pact) SpawnFromInstance(obj actor.Actor, options ...catcher.Options) (*actor.PID, error) {
	props := actor.FromInstance(obj)
	return p.spawn(props, options...)
}

func (p *Pact) SpawnFromProducer(producer actor.Producer, options ...catcher.Options) (*actor.PID, error) {
	props := actor.FromProducer(producer)
	return p.spawn(props, options...)
}

func (p *Pact) SpawnFromFunc(f actor.ActorFunc, options ...catcher.Options) (*actor.PID, error) {
	props := actor.FromFunc(f)
	return p.spawn(props, options...)
}

func (p *Pact) SpawnNullActor(options ...catcher.Options) (*actor.PID, error) {
	return p.SpawnFromInstance(&catcher.NullReceiver{}, options...)
}
