package infra

import (
	"github.com/gusram01/linked-bookmarks/internal/link/domain"
	"gorm.io/gorm"
)

type LinkRepoWithGorm struct {
    db *gorm.DB
}

func NewLinkRepoWithGorm(db *gorm.DB) *LinkRepoWithGorm {
    return &LinkRepoWithGorm{
        db: db,
    }
}

func (lr *LinkRepoWithGorm) Create(r domain.LinkRequest) error{

    return nil
}

func (lr *LinkRepoWithGorm)    GetOneById(id string) (domain.LinkResponse, error){

    return domain.LinkResponse{}, nil
}

func (lr *LinkRepoWithGorm)    GetAll() (domain.LinkResponse, error) {

    return domain.LinkResponse{}, nil
}
