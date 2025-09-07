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
		return domain.Link{}, internal.WrapErrorf(
			err,
			internal.ErrorCodeInvalidField,
			"CreateLink::Invalid::URL::%s",
			r.Url,
		)
	}

	link, err := uc.r.UpsertOne(r)

	if err != nil {
		return domain.Link{}, internal.WrapErrorf(
			err,
			internal.ErrorCodeDBQueryError,
			"CreateLink::Create::Err::ValidateRequest",
		)
	}

	return link, nil
}
