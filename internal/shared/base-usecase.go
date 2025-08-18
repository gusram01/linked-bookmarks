package shared

type QueryBaseUseCase[Request any, Response any] interface {
	Execute(req Request) (Response, error)
}

type CommandBaseUseCase[Request any] interface {
	Execute(req Request) error
}
