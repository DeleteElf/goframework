package entities

// IWorker 工作者接口
type IWorker interface {
	Init(config IConfig) bool
	GetConfig() IConfig
	GetIsRunning() bool
	Start() bool
	Stop() bool
	run()
}

type BaseWorker struct {
	IWorker
	isRunning bool
}

func (worker *BaseWorker) GetIsRunning() bool {
	return worker.isRunning
}

func (worker *BaseWorker) Start() bool {
	worker.isRunning = true
	return true
}

func (worker *BaseWorker) Stop() bool {
	worker.isRunning = false
	return true
}
