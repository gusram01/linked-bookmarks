package usecases

import (
	"github.com/gusram01/linked-bookmarks/internal"
	"github.com/gusram01/linked-bookmarks/internal/link/domain"
)

type GetAllLinks struct {
	r domain.LinkRepository
}

func NewGetAllLinksUse(r domain.LinkRepository) *GetAllLinks {
	return &GetAllLinks{
		r: r,
	}
}

func (uc *GetAllLinks) Execute(r domain.GetAllLinksRequestDto) ([]domain.Link, error) {
	if r.Limit == 0 {
		r.Limit = 5
	}

	links, err := uc.r.GetAll(r)

	if err != nil {
		return []domain.Link{}, internal.WrapErrorf(
			err,
			internal.ErrorCodeDBQueryError,
			"GetAll::Links::Error",
		)
	}

	return links, nil
}
