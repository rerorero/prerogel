package util

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/sirupsen/logrus"
)

// ActorUtil is utility for actor
type ActorUtil struct {
	Logger *logrus.Logger
}

// Fail crashes actor
func (u *ActorUtil) Fail(ctx actor.Context, err error) {
	u.Logger.WithError(err).
		WithField("actor", ctx.Self().Id).
		Error(err.Error())
	// let it crash
	panic(err)
}

// LogError logs error level message
func (u *ActorUtil) LogError(ctx actor.Context, msg string) {
	u.Logger.WithField("actor", ctx.Self().Id).Error(msg)
}

// LogWarn logs warn level message
func (u *ActorUtil) LogWarn(ctx actor.Context, msg string) {
	u.Logger.WithField("actor", ctx.Self().Id).Warn(msg)
}

// LogInfo logs info level message
func (u *ActorUtil) LogInfo(ctx actor.Context, msg string) {
	u.Logger.WithField("actor", ctx.Self().Id).Info(msg)
}

// LogDebug logs debug level message
func (u *ActorUtil) LogDebug(ctx actor.Context, msg string) {
	u.Logger.WithField("actor", ctx.Self().Id).Debug(msg)
}

// IsSystemMessage returns if message is SystemMessage
func (u *ActorUtil) IsSystemMessage(m interface{}) bool {
	if _, ok := m.(actor.AutoReceiveMessage); ok {
		return true
	}
	if _, ok := m.(actor.SystemMessage); ok {
		return true
	}
	return false
}
