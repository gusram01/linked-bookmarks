package usecases

import (
	"github.com/gusram01/linked-bookmarks/internal"
	"github.com/gusram01/linked-bookmarks/internal/link/domain"
)

type CreateOneLink struct {
	r domain.LinkRepository
}

func NewCreateOneLinkUse(r domain.LinkRepository) *CreateOneLink {
	return &CreateOneLink{
		r: r,
	}
}

func (uc *CreateOneLink) Execute(r domain.NewLinkRequestDto) (domain.Link, error) {
	if err := r.Validate(); err != nil {
		return domain.Link{}, internal.NewErrorf(
			internal.ErrorCodeInvalidField,
			"CreateLink::Invalid::URL::%s",
			r.Url,
		)
	}

	// TODO: handle database errors and mapping to internal.ErrorCode
	return uc.r.Create(r)
}
