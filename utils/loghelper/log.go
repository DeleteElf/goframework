package loghelper

import (
	"github.com/pion/logging"
	"log"
	"os"
	"runtime"
	"time"
)

type LogLevel int

const (
	Fatal LogLevel = iota + 1
	Panic
	Error
	Warn
	Info
	Debug
	Trace
)

func (level LogLevel) String() string {
	return [...]string{"Fatal", "Panic", "Error", "Warn", "Info", "Debug", "Trace"}[level-1]
}

func GetLogLevel(level string) LogLevel {
	levelMap := map[string]LogLevel{
		"Fatal": Fatal,
		"Panic": Panic,
		"Error": Error,
		"Warn":  Warn,
		"Info":  Info,
		"Debug": Debug,
		"Trace": Trace,
	}
	return levelMap[level]
}

func GetColorLevel(level string) string {
	levelMap := map[string]string{
		"Fatal": "\x1b[97;31m", // 红色
		"Panic": "\x1b[97;31m", // 红色
		"Error": "\x1b[97;31m", // 红色
		"Warn":  "\x1b[97;33m", // 黄色
		"Info":  "\x1b[97;0m",  // 白色
		"Debug": "\x1b[97;32m", //绿色
	}
	//背景色。。。。
	//levelMap := map[string]string{
	//	"Fatal": "\x1b[97;41m", // 红色
	//	"Panic": "\x1b[97;41m", // 红色
	//	"Error": "\x1b[97;41m", // 红色
	//	"Warn":  "\x1b[97;43m", // 黄色
	//	"Info":  "\x1b[97;0m",  // 白色
	//	"Debug": "\x1b[97;42m", //绿色
	//}
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
	log.Logger
}

func NewLeveledLogger(scope string, level LogLevel) *LeveledLoggerImpl {
	logger := &LeveledLoggerImpl{scope: scope, Level: level}
	logger.SetFlags(0)
	//log.SetFlags(log.Ldate | log.Lmicroseconds)
	logger.SetOutput(os.Stdout)
	//logger.SetPrefix(logManager.levels[level] + "[" + logger.scope + "][" + time.Now().Format("2006-01-02 15:04:05") + "]")
	return logger
}

func NewLogManager(lvl LogLevel) *LogManager {
	return &LogManager{
		Level:      lvl,
		FactoryMap: map[string]*LeveledLoggerImpl{"default": NewLeveledLogger("default", lvl)},
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
	for _, factory := range logM.FactoryMap {
		factory.Level = lvl
	}
}

func (logM *LogManager) NewLogger(scope string) logging.LeveledLogger {
	if logM.FactoryMap[scope] == nil {
		logM.FactoryMap[scope] = NewLeveledLogger(scope, logM.Level)
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

func (logger *LeveledLoggerImpl) SetLevel(level LogLevel) {
	logger.Level = level
}

func (logger *LeveledLoggerImpl) Println(level LogLevel, message string) {
	if logger.Level >= level {
		logger.SetPrefix(GetColorLevel(level.String()) + logManager.levels[level] + "[" + logger.scope + "][" + time.Now().Format("2006-01-02 15:04:05.000") + "]")
		msg := message + "\x1b[0m"
		switch level {
		case Fatal:
			//logger.SetOutput(os.Stderr)
			logger.Fatalln(msg)
		case Panic:
			//logger.SetOutput(os.Stderr)
			logger.Panicln(msg)
		default:
			//logger.SetOutput(os.Stdout)
			logger.Logger.Println(msg)
			break
		}
	}
}
func (logger *LeveledLoggerImpl) Printf(level LogLevel, messageFormat string, args ...any) {
	if logger.Level >= level {
		logger.SetPrefix(GetColorLevel(level.String()) + logManager.levels[level] + "[" + logger.scope + "][" + time.Now().Format("2006-01-02 15:04:05.000") + "]")
		format := messageFormat + "\x1b[0m"
		switch level {
		case Fatal:
			//logger.SetOutput(os.Stderr)
			logger.Logger.Fatalf(format, args...)
		case Panic:
			//logger.SetOutput(os.Stderr)
			logger.Logger.Panicf(format, args...)
		default:
			//logger.SetOutput(os.Stdout)
			logger.Logger.Printf(format, args...)
			break
		}
	}
}

func (logger *LeveledLoggerImpl) Trace(message string) {
	logger.Println(Trace, message)
}
func (logger *LeveledLoggerImpl) Tracef(messageFormat string, args ...any) {
	logger.Printf(Trace, messageFormat, args...)
}

func (logger *LeveledLoggerImpl) Debug(message string) {
	logger.Println(Debug, message)
}
func (logger *LeveledLoggerImpl) Debugf(messageFormat string, args ...any) {
	logger.Printf(Debug, messageFormat, args...)
}

func (logger *LeveledLoggerImpl) Info(message string) {
	logger.Println(Info, message)
}

func (logger *LeveledLoggerImpl) Infof(messageFormat string, args ...any) {
	logger.Printf(Info, messageFormat, args...)
}

func (logger *LeveledLoggerImpl) Warn(message string) {
	logger.Println(Warn, message)
}

func (logger *LeveledLoggerImpl) Warnf(messageFormat string, args ...any) {
	logger.Printf(Warn, messageFormat, args...)
}
func (logger *LeveledLoggerImpl) Error(message string) {
	logger.Println(Error, message)
	log.Println(StackTrace(false))
}

func (logger *LeveledLoggerImpl) Errorf(messageFormat string, args ...any) {
	logger.Printf(Error, messageFormat, args...)
}

func (logger *LeveledLoggerImpl) Fatal(message string) {
	logger.Println(Fatal, message)
	log.Println(StackTrace(false))
}

func (logger *LeveledLoggerImpl) Fatalf(messageFormat string, args ...any) {
	logger.Printf(Fatal, messageFormat, args...)
}

func (logger *LeveledLoggerImpl) Panic(message string) {
	logger.Println(Panic, message)
	log.Println(StackTrace(false))
}

func (logger *LeveledLoggerImpl) Panicf(messageFormat string, args ...any) {
	logger.Printf(Panic, messageFormat, args...)
}
