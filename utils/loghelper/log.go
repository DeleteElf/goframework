package loghelper

import (
	"log"
	"runtime"
)

type LogLevel int

const (
	Fatal LogLevel = iota
	Panic
	Error
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
			Fatal: "[Fatal]",
			Panic: "[Panic]",
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
	log.SetFlags(log.Ldate | log.Lmicroseconds)
}

func StackTrace(all bool) string {
	buf := make([]byte, 10240)
	for {
		size := runtime.Stack(buf, all)
		if size == len(buf) {
			buf = make([]byte, len(buf)<<1)
			continue
		}
		break
	}
	return string(buf)
}

func (logM *LogManager) Println(level LogLevel, message any) {
	//logc.Debugv(log.ctx, message)
	if logM.Level >= level {
		log.SetPrefix(logM.levels[level])
		switch level {
		case Fatal:
			log.Fatalln(message)
		case Panic:
			log.Panicln(message)
		default:
			log.Println(message)
			break
		}
	}
}
func (logM *LogManager) PrintFormat(level LogLevel, messageFormat string, args ...any) {
	//logc.Debugf(log.ctx, messageFormat, args...)
	if logM.Level >= level {
		log.SetPrefix(logM.levels[level])
		log.Printf(messageFormat, args...)
		//message := fmt.Sprintf(messageFormat, args...)
		//log.Println(message)
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
	log.Println(StackTrace(false))
}

func (logM *LogManager) ErrorFormat(messageFormat string, args ...any) {
	//logc.Errorf(log.ctx, messageFormat, args...)
	logM.PrintFormat(Error, messageFormat, args...)
}

func (logM *LogManager) Fatal(message any) {
	//logc.Errorv(log.ctx, message)
	logM.Println(Fatal, message)
	log.Println(StackTrace(false))
}

func (logM *LogManager) FatalFormat(messageFormat string, args ...any) {
	//logc.Errorf(log.ctx, messageFormat, args...)
	logM.PrintFormat(Fatal, messageFormat, args...)
}

func (logM *LogManager) Panic(message any) {
	//logc.Errorv(log.ctx, message)
	logM.Println(Panic, message)
	log.Println(StackTrace(false))
}

func (logM *LogManager) PanicFormat(messageFormat string, args ...any) {
	//logc.Errorf(log.ctx, messageFormat, args...)
	logM.PrintFormat(Panic, messageFormat, args...)
}
