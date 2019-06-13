package util

import (
	"errors"
	"reflect"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/golang/protobuf/proto"
)

type forward struct {
	request     proto.Message
	resExpected proto.Message
	respond     chan proto.Message
}

// ActorProxy is proxy that captures messages sent to Parent() under test
type ActorProxy struct {
	parent     *actor.PID
	underlying *actor.PID
	current    *forward
}

// NewActorProxy returns a new instance
func NewActorProxy(root actor.SpawnerContext, underlyingProps *actor.Props, customHandler func(ctx actor.Context)) *ActorProxy {
	proxy := &ActorProxy{}
	init := make(chan int, 1)
	proxy.parent = root.Spawn(actor.PropsFromFunc(func(ctx actor.Context) {
		switch m := ctx.Message().(type) {
		case *actor.Started:
			proxy.underlying = ctx.Spawn(underlyingProps)
			init <- 1
		case *forward:
			ctx.Request(proxy.underlying, m.request)
			if m.respond != nil && m.resExpected != nil {
				proxy.current = m
			}
		case proto.Message:
			if proxy.current != nil && reflect.TypeOf(m).String() == reflect.TypeOf(proxy.current.resExpected).String() {
				proxy.current.respond <- m
			} else if customHandler != nil {
				customHandler(ctx)
			}
		}
	}))
	<-init
	return proxy
}

// Underlying returns PID under the test
func (proxy *ActorProxy) Underlying() *actor.PID {
	return proxy.underlying
}

// SendAndAwait sends a request to underlying actor and wait for response
func (proxy *ActorProxy) SendAndAwait(ctx actor.SenderContext, req proto.Message, resEmpty proto.Message, timeout time.Duration) (proto.Message, error) {
	msg := &forward{req, resEmpty, make(chan proto.Message, 1)}
	defer close(msg.respond)
	ctx.Request(proxy.parent, msg)
	select {
	case res := <-msg.respond:
		return res, nil
	case <-time.After(timeout):
		return nil, errors.New("sendAndAwait timeout")
	}
}

// Send sends a request via proxy
func (proxy *ActorProxy) Send(ctx actor.SenderContext, req proto.Message) {
	msg := &forward{req, nil, nil}
	ctx.Request(proxy.parent, msg)
}
