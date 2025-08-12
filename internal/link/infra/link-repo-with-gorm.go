package infra

import (
	"github.com/gusram01/linked-bookmarks/internal/link/domain"
	"gorm.io/gorm"
)

type LinkModel struct {
    gorm.Model
    Url string `gorm:"index:idx_link_models_url,unique"`
}

type LinkRepoWithGorm struct {
    db *gorm.DB
}

func NewLinkRepoWithGorm(db *gorm.DB) *LinkRepoWithGorm {
    return &LinkRepoWithGorm{
        db: db,
    }
}

func (lr *LinkRepoWithGorm) Create(r domain.LinkRequest) (domain.Link, error){
    link := LinkModel{Url: string(r.Url)}

    result := lr.db.Create(&link)

    if result.Error != nil {
        return domain.Link{}, result.Error
    }

    return domain.Link{
        ID: link.ID,
        Url: link.Url,
        CreatedAt: link.CreatedAt,
        UpdatedAt: link.UpdatedAt,
        DeletedAt: link.DeletedAt.Time,
    }, nil
}

func (lr *LinkRepoWithGorm)    GetOneById(id uint) (domain.Link, error){

    var link LinkModel

    result := lr.db.First(&link, id)

    if result.Error != nil {
        return domain.Link{}, result.Error
    }

    return domain.Link{
        ID: link.ID,
        Url: link.Url,
        CreatedAt: link.CreatedAt,
        UpdatedAt: link.UpdatedAt,
        DeletedAt: link.DeletedAt.Time,
    }, nil
}

func (lr *LinkRepoWithGorm)    GetAll() ([]domain.Link, error) {

    return []domain.Link{}, nil
}

