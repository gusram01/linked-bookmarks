package domain

type TagRepository interface {
	AddOne(name string) (Tag, error)
	AddMany(r CreateManyTagsRequestDto) ([]Tag, error)
}
