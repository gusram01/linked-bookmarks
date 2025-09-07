package infra

import (
	"github.com/gusram01/linked-bookmarks/internal"
	"github.com/gusram01/linked-bookmarks/internal/shared/models"
	"github.com/gusram01/linked-bookmarks/internal/tag/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TagRepoWithGorm struct {
	db *gorm.DB
}

func NewTagRepoWithGorm(db *gorm.DB) *TagRepoWithGorm {
	return &TagRepoWithGorm{
		db: db,
	}
}

func (tr *TagRepoWithGorm) AddOne(name string) (domain.Tag, error) {
	tag := models.Tag{Name: name}

	result := tr.db.
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "name"}},
			DoUpdates: clause.AssignmentColumns([]string{"updated_at"}),
		}).
		Create(&tag)
	if result.Error != nil {
		return domain.Tag{}, internal.WrapErrorf(result.Error,
			internal.ErrorCodeDBQueryError, "Tag::DB::Create")
	}

	return domain.Tag{
		ID:        tag.ID,
		Name:      tag.Name,
		CreatedAt: tag.CreatedAt,
		UpdatedAt: tag.UpdatedAt,
		DeletedAt: tag.DeletedAt.Time,
	}, nil
}

func (tr *TagRepoWithGorm) AddMany(r domain.CreateManyTagsRequestDto) ([]domain.Tag, error) {
	var tags []models.Tag

	err := tr.db.Transaction(func(tx *gorm.DB) error {
		for _, name := range r.Names {
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
			Where("name IN ?", r.Names).
			Find(&existingTags).Error; err != nil {
			return internal.WrapErrorf(
				err,
				internal.ErrorCodeDBQueryError,
				"Tag::DB::Find::Err::%s",
				err.Error(),
			)
		}

		var link models.Link
		link.ID = r.LinkID

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

	if err != nil {
		return nil, err
	}

	var createdTags []domain.Tag
	for _, tag := range tags {
		createdTags = append(createdTags, domain.Tag{
			ID:        tag.ID,
			Name:      tag.Name,
			CreatedAt: tag.CreatedAt,
			UpdatedAt: tag.UpdatedAt,
			DeletedAt: tag.DeletedAt.Time,
		})
	}

	return createdTags, nil
}
