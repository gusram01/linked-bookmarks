package domain

type LinkRepository interface {
	UpsertOne(r NewLinkRequestDto) (Link, error)
	GetOneById(r GetLinkRequestDto) (Link, error)
	GetAll(r GetAllLinksRequestDto) ([]Link, error)
}
