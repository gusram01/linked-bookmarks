package usecases

import "github.com/gusram01/linked-bookmarks/internal/link/domain"

type GetAllLinks struct {
	r domain.LinkRepository
}

func NewGetAllLinksUse(r domain.LinkRepository) *GetAllLinks {
	return &GetAllLinks{
		r: r,
	}
}

func (uc *GetAllLinks) Execute(cs string) ([]domain.Link, error) {

	return uc.r.GetAll(cs)
}
