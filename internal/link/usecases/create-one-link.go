package usecases

import (
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

func (uc *CreateOneLink) Execute(r domain.LinkRequest) (domain.Link, error) {
    if err := r.Validate(); err != nil {
        return domain.Link{}, err
    }

    return uc.r.Create(r)
}

