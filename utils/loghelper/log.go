package loghelper

import (
	"github.com/pion/logging"
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
	Trace
)

func (level LogLevel) String() string {
	return [...]string{"Fatal", "Panic", "Error", "Warn", "Info", "Debug"}[level]
}

func GetLogLevel(scope string) LogLevel {
	levelMap := map[string]LogLevel{
		"Fatal": Fatal,
		"Panic": Panic,
		"Error": Error,
		"Warn":  Warn,
		"Info":  Info,
		"Debug": Debug,
	}
	return levelMap[scope]
}

type LogManager struct {
	//ctx context.Context
	logging.LoggerFactory
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

func NewLogger(scope string) logging.LeveledLogger {
	lvl := GetLogLevel(scope)
	log.SetFlags(log.Ldate | log.Lmicroseconds)
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
func (logM *LogManager) Printlnf(level LogLevel, messageFormat string, args ...any) {
	if logM.Level >= level {
		log.SetPrefix(logM.levels[level])
		log.Printf(messageFormat, args...)
	}
}

func (logM *LogManager) Trace(message string) {
	logM.Println(Trace, message)
}
func (logM *LogManager) Tracef(messageFormat string, args ...any) {
	logM.Printlnf(Trace, messageFormat, args...)
}

func (logM *LogManager) Debug(message string) {
	logM.Println(Debug, message)
}
func (logM *LogManager) Debugf(messageFormat string, args ...any) {
	logM.Printlnf(Debug, messageFormat, args...)
}

func (logM *LogManager) Info(message string) {
	logM.Println(Info, message)
}

func (logM *LogManager) Infof(messageFormat string, args ...any) {
	logM.Printlnf(Info, messageFormat, args...)
}

func (logM *LogManager) Warn(message string) {
	logM.Println(Warn, message)
}

func (logM *LogManager) Warnf(messageFormat string, args ...any) {
	logM.Printlnf(Warn, messageFormat, args...)
}
func (logM *LogManager) Error(message string) {
	logM.Println(Error, message)
	log.Println(StackTrace(false))
}

func (logM *LogManager) Errorf(messageFormat string, args ...any) {
	logM.Printlnf(Error, messageFormat, args...)
}

func (logM *LogManager) Fatal(message string) {
	logM.Println(Fatal, message)
	log.Println(StackTrace(false))
}

func (logM *LogManager) Fatalf(messageFormat string, args ...any) {
	logM.Printlnf(Fatal, messageFormat, args...)
}

func (logM *LogManager) Panic(message string) {
	logM.Println(Panic, message)
	log.Println(StackTrace(false))
}

func (logM *LogManager) Panicf(messageFormat string, args ...any) {
	logM.Printlnf(Panic, messageFormat, args...)
}
