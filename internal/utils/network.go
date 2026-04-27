package utils

import (
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

// CheckWSOrigin 校验 WebSocket 的 Origin 来源是否安全。
// 默认仅允许同源请求，可通过环境变量 BH_ALLOWED_ORIGINS 配置额外的允许列表（逗号分隔）。
func CheckWSOrigin(r *http.Request) bool {
	origin := r.Header.Get("Origin")
	if origin == "" {
		// 非浏览器发起的请求（如直接用脚本连接）通常不带 Origin，默认放行。
		return true
	}

	u, err := url.Parse(origin)
	if err != nil {
		return false
	}

	// 0. 开发环境校验：如果是非 Release 模式，默认放行
	if gin.Mode() != gin.ReleaseMode {
		return true
	}

	// 1. 同源校验：Origin 的 Host 与请求头中的 Host 一致
	if strings.EqualFold(u.Host, r.Host) {
		return true
	}

	// 2. 环境变量配置的允许列表校验
	allowedOrigins := os.Getenv("BH_ALLOWED_ORIGINS")
	if allowedOrigins != "" {
		origins := strings.Split(allowedOrigins, ",")
		for _, o := range origins {
			o = strings.TrimSpace(o)
			if o == "*" {
				return true
			}
			// 匹配完整 Origin (如 http://localhost:5173) 或仅 Host 部分
			if strings.EqualFold(o, origin) || strings.EqualFold(o, u.Host) {
				return true
			}
		}
	}

	// 3. 兜底策略：如果是开发环境常见的 localhost/127.0.0.1，且端口不一致的情况，
	// 如果用户没有配置允许列表，我们在非 Release 模式下可以考虑放行，
	// 但为了安全，默认应严格限制。建议开发时通过 BH_ALLOWED_ORIGINS=localhost:5173 显式开启。

	return false
}
