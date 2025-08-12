package usecases

type BaseUseCase[Request any, Response any] interface {
    Execute(req Request) (Response, error)
}
