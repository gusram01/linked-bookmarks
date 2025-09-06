package domain

type LinkRepository interface {
	UpsertOne(r NewLinkRequestDto) (Link, error)
	UpdateSummary(id uint, summary string) error
	GetOneById(r GetLinkRequestDto) (Link, error)
	GetAll(r GetAllLinksRequestDto) ([]Link, error)
}
