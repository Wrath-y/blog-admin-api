package def

const (
	// GormDBDataType
	JSON  = "JSON"
	JSONB = "JSONB"

	// layout
	TimeLayout = "2006-01-02 15:04:05"
	DateLayout = "2006-01-02"

	// env
	EnvProduction  = "production"
	EnvDevelopment = "development"
	EnvTesting     = "testing"
	StagingEnv     = "staging"

	CxtUserInfo = "user_info"

	AuthorizationToken  = "Authorization-Token"
	AuthorizationAppKey = "Authorization-AppKey"
	XRequestID          = "X-Request-Id"

	Page     = 1
	PageSize = 100
)
