package worker

type Task interface {
	Process() error
}
