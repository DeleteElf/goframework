package entities

// IWorker 工作者接口
type IWorker interface {
	Init(config IConfig) error
}

// IRunWorker 需要运行的工作者接口
type IRunWorker interface {
	IWorker
	Run()
}
