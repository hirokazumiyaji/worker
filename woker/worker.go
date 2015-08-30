package worker

type Worker interface {
	Start() error
	Shutdown() error
}
