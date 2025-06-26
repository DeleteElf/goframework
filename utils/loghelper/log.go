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

func GetLogLevel(level string) LogLevel {
	levelMap := map[string]LogLevel{
		"Fatal": Fatal,
		"Panic": Panic,
		"Error": Error,
		"Warn":  Warn,
		"Info":  Info,
		"Debug": Debug,
	}
	return levelMap[level]
}

type LogManager struct {
	logging.LoggerFactory
	FactoryMap map[string]*LeveledLoggerImpl
	Level      LogLevel
	levels     map[LogLevel]string
}

type LeveledLoggerImpl struct {
	logging.LeveledLogger
	Level LogLevel
	scope string
	//logr.Logger
}

func NewLogManager(lvl LogLevel) *LogManager {
	return &LogManager{
		Level:      lvl,
		FactoryMap: map[string]*LeveledLoggerImpl{"default": &LeveledLoggerImpl{scope: "default", Level: lvl}},
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
	for _, factory := range logM.FactoryMap {
		factory.Level = lvl
	}
}

func (logM *LogManager) NewLogger(scope string) logging.LeveledLogger {
	if logM.FactoryMap[scope] == nil {
		logM.FactoryMap[scope] = &LeveledLoggerImpl{scope: scope, Level: logM.Level}
	}
	return logM.FactoryMap[scope]
}

func GetDefaultLogger() *LeveledLoggerImpl {
	GetLogManager()
	return logManager.FactoryMap["default"]
}

func GetLogger(scope string) *LeveledLoggerImpl {
	return GetLogManager().NewLogger(scope).(*LeveledLoggerImpl)
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

func (logger *LeveledLoggerImpl) Println(level LogLevel, message any) {
	if logger.Level >= level {
		log.SetPrefix(logManager.levels[level] + "[" + logger.scope + "]")
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
func (logger *LeveledLoggerImpl) Printlnf(level LogLevel, messageFormat string, args ...any) {
	if logger.Level >= level {
		log.SetPrefix(logManager.levels[level] + "[" + logger.scope + "]")
		log.Printf(messageFormat, args...)
	}
}

func (logger *LeveledLoggerImpl) Trace(message string) {
	logger.Println(Trace, message)
}
func (logger *LeveledLoggerImpl) Tracef(messageFormat string, args ...any) {
	logger.Printlnf(Trace, messageFormat, args...)
}

func (logger *LeveledLoggerImpl) Debug(message string) {
	logger.Println(Debug, message)
}
func (logger *LeveledLoggerImpl) Debugf(messageFormat string, args ...any) {
	logger.Printlnf(Debug, messageFormat, args...)
}

func (logger *LeveledLoggerImpl) Info(message string) {
	logger.Println(Info, message)
}

func (logger *LeveledLoggerImpl) Infof(messageFormat string, args ...any) {
	logger.Printlnf(Info, messageFormat, args...)
}

func (logger *LeveledLoggerImpl) Warn(message string) {
	logger.Println(Warn, message)
}

func (logger *LeveledLoggerImpl) Warnf(messageFormat string, args ...any) {
	logger.Printlnf(Warn, messageFormat, args...)
}
func (logger *LeveledLoggerImpl) Error(message string) {
	logger.Println(Error, message)
	log.Println(StackTrace(false))
}

func (logger *LeveledLoggerImpl) Errorf(messageFormat string, args ...any) {
	logger.Printlnf(Error, messageFormat, args...)
}

func (logger *LeveledLoggerImpl) Fatal(message string) {
	logger.Println(Fatal, message)
	log.Println(StackTrace(false))
}

func (logger *LeveledLoggerImpl) Fatalf(messageFormat string, args ...any) {
	logger.Printlnf(Fatal, messageFormat, args...)
}

func (logger *LeveledLoggerImpl) Panic(message string) {
	logger.Println(Panic, message)
	log.Println(StackTrace(false))
}

func (logger *LeveledLoggerImpl) Panicf(messageFormat string, args ...any) {
	logger.Printlnf(Panic, messageFormat, args...)
}
