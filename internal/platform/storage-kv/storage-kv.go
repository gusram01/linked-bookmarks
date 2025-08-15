package storagekv

import (
	"github.com/gofiber/storage/cloudflarekv"
	"github.com/gusram01/linked-bookmarks/internal/platform/config"
)

var kvStorage *cloudflarekv.Storage

func GetStorage() *cloudflarekv.Storage {
	if kvStorage == nil {
		kvStorage = cloudflarekv.New(cloudflarekv.Config{
			Key:         config.Config("GC_MARK_KV_STORAGE_TOKEN"),
			AccountID:   config.Config("GC_MARK_CF_ACCOUNT_ID"),
			NamespaceID: config.Config("GC_MARK_CF_NAMESPACE_ID"),
			Email:       config.Config("GC_MARK_CF_EMAIL"),
			Reset:       false,
		})

		return kvStorage
	}

	return kvStorage
}
