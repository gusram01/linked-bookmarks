package infra

import (
	"strings"

	"github.com/gusram01/linked-bookmarks/internal"
	"github.com/gusram01/linked-bookmarks/internal/onboarding/domain"
	"github.com/gusram01/linked-bookmarks/internal/onboarding/infra/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserRepoWithGorm struct {
	db *gorm.DB
}

func NewUserRepoWithGorm(db *gorm.DB) *UserRepoWithGorm {
	return &UserRepoWithGorm{
		db: db,
	}
}

func (ur *UserRepoWithGorm) Upsert(u *domain.User) error {
	user := models.User{
		AuthID: u.AuthID,
	}

	if err := ur.db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&user).Error; err != nil {

		if strings.Contains(
			err.Error(),
			"violates unique constraint \"idx_user_auth_id\"",
		) {
			return internal.WrapErrorf(
				err,
				internal.ErrorCodeWHHandleUserFound,
				"The User Auth ID already exists: %s",
				user.AuthID,
			)
		}

		return err
	}

	u.ID = user.ID
	u.CreatedAt = user.CreatedAt
	u.UpdatedAt = user.UpdatedAt
	u.DeletedAt = user.DeletedAt.Time

	return nil
}
