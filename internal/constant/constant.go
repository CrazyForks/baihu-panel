package constant

const (

	// ConfigPath 配置文件路径
	ConfigPath = "configs/config.json"

	// DataDir 数据目录
	DataDir = "./data"

	// WebDistDir 前端构建目录
	WebDistDir = "./web/dist"

	// DefaultRole 默认用户角色
	DefaultRole = "user"

	// AdminRole 管理员角色
	AdminRole = "admin"

	// DefaultTablePrefix 默认表前缀
	DefaultTablePrefix = "baihu_"

	// ScriptsWorkDir 脚本工作目录
	ScriptsWorkDir = "./data/scripts"

	// DefaultPageSize 默认分页大小
	DefaultPageSize = 10

	// CookieName Cookie 名称
	CookieName = "BHToken"

	// TokenExpireDays Token 过期天数
	TokenExpireDays = 7
	// CookieMaxAge Cookie 有效期（秒）7天
	CookieMaxAge = 86400 * TokenExpireDays

	// DefaultJWTSecret 默认 JWT 密钥
	DefaultJWTSecret = "baihu-default-secret-key"

	// DefaultTaskTimeout 默认任务超时时间（分钟）
	DefaultTaskTimeout = 30
)

// TablePrefix 表前缀，可在运行时设置
var TablePrefix = DefaultTablePrefix

// JWTSecret JWT 密钥，可通过配置文件设置
var JWTSecret = DefaultJWTSecret
