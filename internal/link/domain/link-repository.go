package domain

type LinkRepository interface {
	UpsertOne(r NewLinkRequestDto) (Link, error)
	UpdateSummary(r UpdateSummaryRequestDto) error
	UpdateTags(r UpdateTagsRequestDto) error
	GetOneById(r GetLinkRequestDto) (Link, error)
	GetAll(r GetAllLinksRequestDto) (GetAllQueryResultDto, error)
	GetManyByIds(r GetManyLinksByIdsRequestDto) ([]Link, error)
}
