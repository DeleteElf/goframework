package loghelper

import (
	"fmt"
	"log"
)

type LogLevel int

const (
	Error LogLevel = iota
	Warn
	Info
	Debug
)

type LogManager struct {
	//ctx context.Context
	Level  LogLevel
	levels map[LogLevel]string
}

func NewLogManager(lvl LogLevel) *LogManager {
	return &LogManager{
		Level: lvl,
		levels: map[LogLevel]string{
			Error: "[Error]",
			Warn:  "[Warn]",
			Info:  "[Info]",
			Debug: "[Debug]",
		},
	}
}

var logManager *LogManager

func GetLogManager() *LogManager {
	if logManager == nil {
		logManager = NewLogManager(Warn)
	}
	return logManager
}

func (logM *LogManager) Init(lvl LogLevel) {
	logM.Level = lvl
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lshortfile)
}

func (logM *LogManager) Println(level LogLevel, message any) {
	//logc.Debugv(log.ctx, message)
	if logM.Level >= level {
		log.SetPrefix(logM.levels[level])
		log.Println(message)
	}
}
func (logM *LogManager) PrintFormat(level LogLevel, messageFormat string, args ...any) {
	//logc.Debugf(log.ctx, messageFormat, args...)
	if logM.Level >= level {
		log.SetPrefix(logM.levels[level])
		message := fmt.Sprintf(messageFormat, args...)
		log.Println(message)
	}
}

func (logM *LogManager) Debug(message any) {
	//logc.Debugv(log.ctx, message)
	//if logM.Level >= Debug {
	//	log.Println(message)
	//}
	logM.Println(Debug, message)
}
func (logM *LogManager) DebugFormat(messageFormat string, args ...any) {
	//logc.Debugf(log.ctx, messageFormat, args...)
	//if logM.Level >= Debug {
	//	log.Printf(messageFormat, args...)
	//	log.Println()
	//}
	logM.PrintFormat(Debug, messageFormat, args...)
}

func (logM *LogManager) Info(message any) {
	//logc.Infov(log.ctx, message)
	logM.Println(Info, message)
}

func (logM *LogManager) InfoFormat(messageFormat string, args ...any) {
	//logc.Infof(log.ctx, messageFormat, args...)
	logM.PrintFormat(Info, messageFormat, args...)
}

func (logM *LogManager) Warn(message any) {
	//logc.Infov(log.ctx, message)
	logM.Println(Warn, message)
}

func (logM *LogManager) WarnFormat(messageFormat string, args ...any) {
	//logc.Infof(log.ctx, messageFormat, args...)
	logM.PrintFormat(Warn, messageFormat, args...)
}
func (logM *LogManager) Error(message any) {
	//logc.Errorv(log.ctx, message)
	logM.Println(Error, message)
}

func (logM *LogManager) ErrorFormat(messageFormat string, args ...any) {
	//logc.Errorf(log.ctx, messageFormat, args...)
	logM.PrintFormat(Error, messageFormat, args...)
}
