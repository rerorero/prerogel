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
func (u *ActorUtil) Fail(err error) {
	u.Logger.WithError(err).Error(err.Error())
	// let it crash
	panic(err)
}

// LogError logs error level message
func (u *ActorUtil) LogError(msg string) {
	u.Logger.Error(msg)
}

// LogWarn logs warn level message
func (u *ActorUtil) LogWarn(msg string) {
	u.Logger.Warn(msg)
}

// LogInfo logs info level message
func (u *ActorUtil) LogInfo(msg string) {
	u.Logger.Info(msg)
}

// LogDebug logs debug level message
func (u *ActorUtil) LogDebug(msg string) {
	u.Logger.Debug(msg)
}

// AppendLoggerField appends a field to logger
func (u *ActorUtil) AppendLoggerField(key string, value interface{}) {
	u.Logger = u.Logger.WithField(key, value).Logger
}

// IsSystemMessage returns if message is SystemMessage
func (u *ActorUtil) IsSystemMessage(m interface{}) bool {
	_, ok := m.(actor.SystemMessage)
	return ok
}
