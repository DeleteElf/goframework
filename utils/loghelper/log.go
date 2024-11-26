package loghelper

import (
	"context"
	"github.com/zeromicro/go-zero/core/logc"
)

type LogManager struct {
	ctx context.Context
}

func NewLogManager(ctx context.Context) *LogManager {
	return &LogManager{
		ctx: ctx,
	}
}

var logManager *LogManager

func GetLogManager() *LogManager {
	if logManager == nil {
		logManager = NewLogManager(nil)
	}
	return logManager
}

func (log *LogManager) Init(ctx context.Context) {
	log.ctx = ctx
}

func (log *LogManager) Debug(message any) {
	logc.Debugv(log.ctx, message)
}
func (log *LogManager) DebugFormat(messageFormat string, args ...any) {
	logc.Debugf(log.ctx, messageFormat, args...)
}

func (log *LogManager) Info(message any) {
	logc.Infov(log.ctx, message)
}

func (log *LogManager) InfoFormat(messageFormat string, args ...any) {
	logc.Infof(log.ctx, messageFormat, args...)
}
func (log *LogManager) Error(message any) {
	logc.Errorv(log.ctx, message)
}

func (log *LogManager) ErrorFormat(messageFormat string, args ...any) {
	logc.Errorf(log.ctx, messageFormat, args...)
}
