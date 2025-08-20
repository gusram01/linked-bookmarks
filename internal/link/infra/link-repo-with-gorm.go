package infra

import (
	"github.com/gusram01/linked-bookmarks/internal"
	"github.com/gusram01/linked-bookmarks/internal/link/domain"
	"github.com/gusram01/linked-bookmarks/internal/shared/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
	var user models.User
	var link models.Link

	user = models.User{AuthID: r.Subject}

	userResult := lr.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "auth_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"updated_at"}),
	}).Create(&user)

	if userResult.Error != nil {
		return domain.Link{}, internal.WrapErrorf(userResult.Error,
			internal.ErrorCodeDBQueryError, "User::DB::Upsert")
	}

	link = models.Link{Url: string(r.Url)}

	linkResult := lr.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "url"}},
		DoUpdates: clause.AssignmentColumns([]string{"updated_at"}),
	}).Create(&link)

	if linkResult.Error != nil {
		return domain.Link{}, internal.WrapErrorf(linkResult.Error,
			internal.ErrorCodeDBQueryError, "Link::DB::Upsert")
	}

	associationResult := lr.db.Model(&link).Association("Users").Append(&user)

	if associationResult != nil {
		return domain.Link{}, internal.WrapErrorf(associationResult,
			internal.ErrorCodeDBQueryError, "Association::DB::Create")
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

func (lr *LinkRepoWithGorm) GetAll(r domain.GetAllLinksRequestDto) ([]domain.Link, error) {

	var ls []domain.Link

	qAllRes := lr.db.
		Model(&models.Link{}).
		Offset(int(r.Offset)).
		Limit(int(r.Limit)).
		Select("links.id, links.url").
		Joins("JOIN user_links ul ON links.id = ul.link_id").
		Joins("JOIN users u ON ul.user_id = u.id ").
		Where("u.auth_id = ?", r.Subject).
		Scan(&ls)

	if qAllRes.Error != nil {
		return []domain.Link{}, qAllRes.Error
	}

	return ls, nil
}
