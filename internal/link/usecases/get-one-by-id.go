package usecases

import (
	"github.com/gusram01/linked-bookmarks/internal"
	"github.com/gusram01/linked-bookmarks/internal/link/domain"
)

type GetOneByIdLink struct {
    r domain.LinkRepository
}

func NewGetOneByIdLinkUse(r domain.LinkRepository) *GetOneByIdLink {
    return &GetOneByIdLink{
        r: r,
    }
}

func (uc *GetOneByIdLink) Execute(id uint) (domain.Link, error) {
    if id <= 0 {
        return domain.Link{}, internal.NewErrorf(
            internal.ErrorCodeInvalidField,
            "GetOneByIdLink::Invalid::ID::%d",
            id,
        )
    }

    // TODO: handle database errors and mapping to internal.ErrorCode
    return uc.r.GetOneById(id)
}
