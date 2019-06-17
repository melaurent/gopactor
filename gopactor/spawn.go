package gopactor

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/melaurent/gopactor/catcher"
	"github.com/melaurent/gopactor/options"
)

func (p *Gopactor) spawn(props *actor.Props, opts ...options.Options) *actor.PID {
	catcher := catcher.New()
	ctx := actor.EmptyRootContext
	pid := catcher.Spawn(ctx, props, opts...)

	p.CatchersByPID[pid.String()] = catcher

	return pid
}

func (p *Gopactor) SpawnFromProducer(producer actor.Producer, opts ...options.Options) *actor.PID {
	props := actor.PropsFromProducer(producer)
	return p.spawn(props, opts...)
}

func (p *Gopactor) SpawnFromFunc(f actor.ActorFunc, opts ...options.Options) *actor.PID {
	props := actor.PropsFromFunc(f)
	return p.spawn(props, opts...)
}

func (p *Gopactor) SpawnNullActor(opts ...options.Options) *actor.PID {
	return p.SpawnFromProducer(func() actor.Actor { return &catcher.NullReceiver{} }, opts...)
}
