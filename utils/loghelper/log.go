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

func (log *LogManager) Debug(messageFormat string, args ...any) {
	logc.Debugf(log.ctx, messageFormat, args)
}

func (log *LogManager) Info(messageFormat string, args ...any) {
	logc.Infof(log.ctx, messageFormat, args)
}

func (log *LogManager) Error(messageFormat string, args ...any) {
	logc.Errorf(log.ctx, messageFormat, args)
}

func (log *LogManager) ErrorV(err error) {
	logc.Errorv(log.ctx, err)
}
