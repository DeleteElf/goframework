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
	//ctx context.Context
	FactoryMap map[string]*LoggerFactoryImpl
	Level      LogLevel
	levels     map[LogLevel]string
}

type LoggerFactoryImpl struct {
	logging.LoggerFactory
	Level LogLevel
}

func NewLogManager(lvl LogLevel) *LogManager {
	return &LogManager{
		Level:      lvl,
		FactoryMap: map[string]*LoggerFactoryImpl{"default": &LoggerFactoryImpl{}},
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
		logM.FactoryMap[scope] = &LoggerFactoryImpl{}
	}
	return logM.FactoryMap[scope]
}

func GetDefaultLogger() *LoggerFactoryImpl {
	GetLogManager()
	return logManager.FactoryMap["default"]
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

func (logger *LoggerFactoryImpl) Println(level LogLevel, message any) {
	if logger.Level >= level {
		log.SetPrefix(logManager.levels[level])
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
func (logger *LoggerFactoryImpl) Printlnf(level LogLevel, messageFormat string, args ...any) {
	if logger.Level >= level {
		log.SetPrefix(logManager.levels[level])
		log.Printf(messageFormat, args...)
	}
}

func (logger *LoggerFactoryImpl) Trace(message string) {
	logger.Println(Trace, message)
}
func (logger *LoggerFactoryImpl) Tracef(messageFormat string, args ...any) {
	logger.Printlnf(Trace, messageFormat, args...)
}

func (logger *LoggerFactoryImpl) Debug(message string) {
	logger.Println(Debug, message)
}
func (logger *LoggerFactoryImpl) Debugf(messageFormat string, args ...any) {
	logger.Printlnf(Debug, messageFormat, args...)
}

func (logger *LoggerFactoryImpl) Info(message string) {
	logger.Println(Info, message)
}

func (logger *LoggerFactoryImpl) Infof(messageFormat string, args ...any) {
	logger.Printlnf(Info, messageFormat, args...)
}

func (logger *LoggerFactoryImpl) Warn(message string) {
	logger.Println(Warn, message)
}

func (logger *LoggerFactoryImpl) Warnf(messageFormat string, args ...any) {
	logger.Printlnf(Warn, messageFormat, args...)
}
func (logger *LoggerFactoryImpl) Error(message string) {
	logger.Println(Error, message)
	log.Println(StackTrace(false))
}

func (logger *LoggerFactoryImpl) Errorf(messageFormat string, args ...any) {
	logger.Printlnf(Error, messageFormat, args...)
}

func (logger *LoggerFactoryImpl) Fatal(message string) {
	logger.Println(Fatal, message)
	log.Println(StackTrace(false))
}

func (logger *LoggerFactoryImpl) Fatalf(messageFormat string, args ...any) {
	logger.Printlnf(Fatal, messageFormat, args...)
}

func (logger *LoggerFactoryImpl) Panic(message string) {
	logger.Println(Panic, message)
	log.Println(StackTrace(false))
}

func (logger *LoggerFactoryImpl) Panicf(messageFormat string, args ...any) {
	logger.Printlnf(Panic, messageFormat, args...)
}
