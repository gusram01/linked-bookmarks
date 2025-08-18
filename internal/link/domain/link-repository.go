package domain

type LinkRepository interface {
	Create(r NewLinkRequestDto) (Link, error)
	GetOneById(r GetLinkRequestDto) (Link, error)
	GetAll(cs string) ([]Link, error)
}
