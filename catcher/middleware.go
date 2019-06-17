package catcher

import (
	"github.com/AsynkronIT/protoactor-go/actor"
)

func (catcher *Catcher) contextDecorator(next actor.ContextDecoratorFunc) actor.ContextDecoratorFunc {
	return func(ctx actor.Context) actor.Context {
		catcherCtx := NewContext(catcher, ctx)
		return next(catcherCtx)
	}
}

func (catcher *Catcher) inboundMiddleware(next actor.ReceiverFunc) actor.ReceiverFunc {
	return func(ctx actor.ReceiverContext, envelope *actor.MessageEnvelope) {
		catcher.processInboundMessage(ctx, envelope)
		next(ctx, envelope)
	}
}

func (catcher *Catcher) processInboundMessage(ctx actor.ReceiverContext, mEnvelope *actor.MessageEnvelope) {
	envelope := &Envelope{
		Sender:  mEnvelope.Sender,
		Target:  ctx.Self(),
		Message: mEnvelope.Message,
	}

	if !isSystemMessage(mEnvelope.Message) {
		if catcher.Options.InboundInterceptionEnabled {
			catcher.ChUserInbound <- envelope
		}
	} else {
		if catcher.Options.SystemInterceptionEnabled {
			catcher.processSystemMessage(envelope)
		}
	}
}

func (catcher *Catcher) processSystemMessage(envelope *Envelope) {
	catcher.ChSystemInbound <- envelope
}

func (catcher *Catcher) outboundMiddleware(next actor.SenderFunc) actor.SenderFunc {
	return func(ctx actor.SenderContext, target *actor.PID, env *actor.MessageEnvelope) {
		catcher.processOutboundMessage(ctx, target, env)
		next(ctx, target, env)
	}
}

func (catcher *Catcher) processOutboundMessage(ctx actor.SenderContext, target *actor.PID, env *actor.MessageEnvelope) {
	// TODO: Is there a difference between using ctx.Message() and env.Message?
	message := env.Message

	if !isSystemMessage(message) {
		catcher.ChUserOutbound <- &Envelope{
			Sender:  ctx.Self(),
			Target:  target,
			Message: message,
		}
	}
}

func isSystemMessage(msg interface{}) bool {
	switch msg.(type) {
	case actor.AutoReceiveMessage:
		return true
	case actor.SystemMessage:
		return true
	}

	return false
}
