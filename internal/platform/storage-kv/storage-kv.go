package storagekv

import (
	"github.com/gofiber/storage/cloudflarekv"
	"github.com/gusram01/linked-bookmarks/internal/platform/config"
)

var kvStorage *cloudflarekv.Storage

func GetStorage() *cloudflarekv.Storage {
	if kvStorage == nil {
		kvStorage = cloudflarekv.New(cloudflarekv.Config{
			Key:         config.ENVS.KvStorageToken,
			AccountID:   config.ENVS.CfAccountId,
			NamespaceID: config.ENVS.CfNamespaceId,
			Email:       config.ENVS.CfEmail,
			Reset:       false,
		})

		return kvStorage
	}

	return kvStorage
}
