package controllers

import (
	"baihu/internal/constant"
	"baihu/internal/database"
	"baihu/internal/models"
	"baihu/internal/services"
	"baihu/internal/utils"
	"fmt"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
)

type SettingsController struct {
	userService *services.UserService
}

func NewSettingsController(userService *services.UserService) *SettingsController {
	return &SettingsController{userService: userService}
}

// ChangePassword 修改密码
func (sc *SettingsController) ChangePassword(c *gin.Context) {
	var req struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	// 暂时使用固定用户名 admin
	user := sc.userService.GetUserByUsername("admin")
	if user == nil {
		utils.NotFound(c, "用户不存在")
		return
	}

	if !sc.userService.ValidatePassword(user, req.OldPassword) {
		utils.BadRequest(c, "原密码错误")
		return
	}

	if err := sc.userService.UpdatePassword(user.ID, req.NewPassword); err != nil {
		utils.ServerError(c, "修改密码失败")
		return
	}

	utils.SuccessMsg(c, "密码修改成功")
}

// CleanLogs 清理日志
func (sc *SettingsController) CleanLogs(c *gin.Context) {
	var req struct {
		Days int `json:"days" binding:"required,min=1"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	cutoff := time.Now().AddDate(0, 0, -req.Days)
	result := database.DB.Where("created_at < ?", cutoff).Delete(&models.TaskLog{})

	utils.Success(c, gin.H{
		"deleted": result.RowsAffected,
	})
}

// GetSiteSettings 获取站点设置
func (sc *SettingsController) GetSiteSettings(c *gin.Context) {
	config := services.Config
	if config == nil {
		utils.Success(c, gin.H{
			"site_name": "白虎面板",
			"port":      8080,
		})
		return
	}

	utils.Success(c, gin.H{
		"site_name": config.Server.SiteName,
		"port":      config.Server.Port,
	})
}

// GetAbout 获取关于信息
func (sc *SettingsController) GetAbout(c *gin.Context) {
	var taskCount, logCount, envCount int64
	database.DB.Model(&models.Task{}).Count(&taskCount)
	database.DB.Model(&models.TaskLog{}).Count(&logCount)
	database.DB.Model(&models.EnvironmentVariable{}).Count(&envCount)

	// 内存使用
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	memUsage := formatBytes(m.Alloc)

	// 运行时间
	uptime := formatDuration(time.Since(constant.StartTime))

	utils.Success(c, gin.H{
		"version":    constant.Version,
		"build_time": constant.BuildTime,
		"mem_usage":  memUsage,
		"uptime":     uptime,
		"task_count": taskCount,
		"log_count":  logCount,
		"env_count":  envCount,
	})
}

// formatBytes 格式化字节数
func formatBytes(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := uint64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// formatDuration 格式化时间间隔
func formatDuration(d time.Duration) string {
	days := int(d.Hours()) / 24
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60

	if days > 0 {
		return fmt.Sprintf("%d天%d小时%d分钟%d秒", days, hours, minutes, seconds)
	}
	if hours > 0 {
		return fmt.Sprintf("%d小时%d分钟%d秒", hours, minutes, seconds)
	}
	if minutes > 0 {
		return fmt.Sprintf("%d分钟%d秒", minutes, seconds)
	}
	return fmt.Sprintf("%d秒", seconds)
}
