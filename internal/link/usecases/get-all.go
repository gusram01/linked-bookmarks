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

func (uc *GetAllLinks) Execute(r domain.GetAllLinksRequestDto) (domain.GetAllQueryResultDto, error) {
	qr, err := uc.r.GetAll(r)

	if err != nil {
		return domain.GetAllQueryResultDto{}, internal.WrapErrorf(
			err,
			internal.ErrorCodeDBQueryError,
			"GetAll::Links::Error",
		)
	}

	pages := int64(qr.TotalCount) / int64(r.Limit)

	if int64(qr.TotalCount)%int64(r.Limit) != 0 {
		pages += 1
	}

	qr.Pages = pages

	return qr, nil
}
