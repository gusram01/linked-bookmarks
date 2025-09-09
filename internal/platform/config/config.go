package config

import (
	"github.com/gusram01/linked-bookmarks/internal"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type EnvVarConfig struct {
	ApiPort                string `envconfig:"GC_MARK_PORT" default:"4200"`
	LogLevel               string `envconfig:"GC_MARK_LOG_LEVEL" default:"0"`
	DbUser                 string `envconfig:"GC_MARK_DB_USER" required:"true"`
	DbPass                 string `envconfig:"GC_MARK_DB_PASS" required:"true"`
	DbName                 string `envconfig:"GC_MARK_DB_NAME" required:"true"`
	DbHost                 string `envconfig:"GC_MARK_DB_HOST" required:"true"`
	DbSSLMode              string `envconfig:"GC_MARK_DB_SSL_MODE" default:"disable"`
	AuthProviderApiKey     string `envconfig:"GC_MARK_AUTH_KEY" required:"true"`
	KvStorageToken         string `envconfig:"GC_MARK_KV_STORAGE_TOKEN" required:"true"`
	CfAccountId            string `envconfig:"GC_MARK_CF_ACCOUNT_ID" required:"true"`
	CfNamespaceId          string `envconfig:"GC_MARK_CF_NAMESPACE_ID" required:"true"`
	CfEmail                string `envconfig:"GC_MARK_CF_EMAIL" required:"true"`
	SentryDsn              string `envconfig:"GC_MARK_SENTRY_DSN" required:"true"`
	WebhookProviderSecret  string `envconfig:"GC_MARK_CLERK_WH_SIGNING_SECRET" required:"true"`
	GeminiApiKey           string `envconfig:"GC_MARK_GEMINI_API_KEY" required:"true"`
	VectorDBHost           string `envconfig:"GC_MARK_VECTOR_DB_HOST" default:"http://localhost"`
	VectorDBPort           string `envconfig:"GC_MARK_VECTOR_DB_PORT" default:"8000"`
	VectorDBCollectionName string `envconfig:"GC_MARK_VECTOR_DB_COLLECTION_NAME" default:"linked-bookmarks"`
}

var ENVS EnvVarConfig

func LoadConfigFile(filename ...string) error {

	if err := godotenv.Load(filename...); err != nil {
		return internal.NewErrorf(internal.ErrorCodeInternal, "failed to load env file:", err)
	}

	if err := envconfig.Process("GC_BOOKMARK_API", &ENVS); err != nil {
		return internal.NewErrorf(internal.ErrorCodeInternal, "failed to process env vars:", err)
	}

	return nil
}
