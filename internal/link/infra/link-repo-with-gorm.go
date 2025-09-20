package infra

import (
	"github.com/gusram01/linked-bookmarks/internal"
	"github.com/gusram01/linked-bookmarks/internal/link/domain"
	"github.com/gusram01/linked-bookmarks/internal/platform/logger"
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

func (lr *LinkRepoWithGorm) UpsertOne(r domain.NewLinkRequestDto) (domain.Link, error) {

	user := models.User{AuthID: r.Subject}

	userResult := lr.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "auth_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"updated_at"}),
	}).Create(&user)

	if userResult.Error != nil {
		return domain.Link{}, internal.WrapErrorf(userResult.Error,
			internal.ErrorCodeDBQueryError, "User::DB::Upsert")
	}

	link := models.Link{Url: string(r.Url)}

	linkResult := lr.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "url"}},
		DoUpdates: clause.Assignments(map[string]interface{}{"attempts": gorm.Expr("links.attempts + ?", 1), "updated_at": gorm.Expr("CURRENT_TIMESTAMP")}),
	}).Create(&link)

	logger.GetLogger().Info("UpsertOne link", "linkID", link.ID, "linkURL", link.Url, "linkAttempts", link.Attempts, "linkSummary", link.Summary)

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

func (lr *LinkRepoWithGorm) UpdateSummary(r domain.UpdateSummaryRequestDto) error {
	result := lr.db.Model(&models.Link{}).Where("id = ?", r.ID).Update("summary", r.Summary)

	if result.Error != nil {
		return internal.WrapErrorf(result.Error,
			internal.ErrorCodeDBQueryError, "Link::DB::UpdateSummary")
	}

	return nil
}

func (lr *LinkRepoWithGorm) UpdateTags(r domain.UpdateTagsRequestDto) error {

	var tags []models.Tag

	return lr.db.
		Transaction(func(tx *gorm.DB) error {

			for _, name := range r.Tags {
				tag := models.Tag{Name: name}
				tags = append(tags, tag)
			}

			tagResult := tx.
				Clauses(clause.OnConflict{
					Columns:   []clause.Column{{Name: "name"}},
					DoUpdates: clause.AssignmentColumns([]string{"updated_at"}),
				}).
				Create(&tags)

			if tagResult.Error != nil {
				return internal.WrapErrorf(tagResult.Error,
					internal.ErrorCodeDBQueryError, "Tag::DB::Create")
			}

			var existingTags []models.Tag

			if err := tx.
				Where("name IN ?", r.Tags).
				Find(&existingTags).Error; err != nil {
				return internal.WrapErrorf(
					err,
					internal.ErrorCodeDBQueryError,
					"Tag::DB::Find::Err::%s",
					err.Error(),
				)
			}

			var link models.Link
			link.ID = r.ID

			appendResult := tx.Model(&link).Association("Tags").Append(existingTags)

			if appendResult != nil {
				return internal.WrapErrorf(
					appendResult,
					internal.ErrorCodeDBQueryError,
					"Tag::DB::Append::Err::%s",
					appendResult.Error(),
				)
			}

			return nil
		})

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
		Summary:   link.Summary,
		Attempts:  link.Attempts,
		CreatedAt: link.CreatedAt,
		UpdatedAt: link.UpdatedAt,
		DeletedAt: link.DeletedAt.Time,
	}, nil
}

func (lr *LinkRepoWithGorm) GetAll(r domain.GetAllLinksRequestDto) (domain.GetAllQueryResultDto, error) {

	var qr domain.GetAllQueryResultDto

	qCountRes := lr.db.
		Model(&models.Link{}).
		Joins("JOIN user_links ul ON links.id = ul.link_id").
		Joins("JOIN users u ON ul.user_id = u.id ").
		Where("u.auth_id = ?", r.Subject).
		Count(&qr.TotalCount)

	if qCountRes.Error != nil {
		return domain.GetAllQueryResultDto{}, qCountRes.Error
	}

	qAllRes := lr.db.
		Model(&models.Link{}).
		Offset(int(r.Offset)).
		Limit(int(r.Limit)).
		Select("links.id, links.url, links.summary, links.attempts").
		Joins("JOIN user_links ul ON links.id = ul.link_id").
		Joins("JOIN users u ON ul.user_id = u.id ").
		Where("u.auth_id = ?", r.Subject).
		Scan(&qr.Links)

	if qAllRes.Error != nil {
		return domain.GetAllQueryResultDto{}, qAllRes.Error
	}

	return qr, nil
}

func (lr *LinkRepoWithGorm) GetManyByIds(r domain.GetManyLinksByIdsRequestDto) ([]domain.Link, error) {

	var links []models.Link

	result := lr.db.
		Model(&models.Link{}).
		Select("links.id, links.url, links.summary, links.attempts").
		Joins("JOIN user_links ul ON links.id = ul.link_id").
		Joins("JOIN users u ON ul.user_id = u.id ").
		Where("u.auth_id = ? AND links.id IN ?", r.Subject, r.IDs).
		Scan(&links)

	if result.Error != nil {
		return []domain.Link{}, result.Error
	}

	var domainLinks []domain.Link

	for _, link := range links {
		domainLinks = append(domainLinks, domain.Link{
			ID:        link.ID,
			Url:       link.Url,
			Summary:   link.Summary,
			Attempts:  link.Attempts,
			CreatedAt: link.CreatedAt,
			UpdatedAt: link.UpdatedAt,
			DeletedAt: link.DeletedAt.Time,
		})
	}

	return domainLinks, nil
}
