package infra

import (
	"github.com/gusram01/linked-bookmarks/internal/link/domain"
	"github.com/gusram01/linked-bookmarks/internal/link/infra/models"
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

func (lr *LinkRepoWithGorm) Create(r domain.NewLinkRequestDto) (domain.Link, error) {
	link := models.Link{Url: string(r.Url)}

	result := lr.db.Create(&link)

	if result.Error != nil {
		return domain.Link{}, result.Error
	}

	return domain.Link{
		ID:        link.ID,
		Url:       link.Url,
		CreatedAt: link.CreatedAt,
		UpdatedAt: link.UpdatedAt,
		DeletedAt: link.DeletedAt.Time,
	}, nil
}

func (lr *LinkRepoWithGorm) GetOneById(r domain.GetLinkRequestDto) (domain.Link, error) {

	var link models.Link

	result := lr.db.First(&link, r.ID)

	if result.Error != nil {
		return domain.Link{}, result.Error
	}

	return domain.Link{
		ID:        link.ID,
		Url:       link.Url,
		CreatedAt: link.CreatedAt,
		UpdatedAt: link.UpdatedAt,
		DeletedAt: link.DeletedAt.Time,
	}, nil
}

func (lr *LinkRepoWithGorm) GetAll(cs string) ([]domain.Link, error) {

	var links []models.Link

	result := lr.db.Where("deleted_at IS null").Find(&links)

	if result.Error != nil {
		return []domain.Link{}, result.Error
	}

	var domainLinks []domain.Link

	for _, link := range links {
		domainLinks = append(domainLinks, domain.Link{
			ID:  link.ID,
			Url: link.Url,
		})
	}

	return domainLinks, nil
}
