package controllers

import (
	"baihu/internal/constant"
	"baihu/internal/database"
	"baihu/internal/models"
	"baihu/internal/services"
	"baihu/internal/utils"

	"github.com/gin-gonic/gin"
)

type DashboardController struct {
	cronService     *services.CronService
	executorService *services.ExecutorService
}

func NewDashboardController(cronService *services.CronService, executorService *services.ExecutorService) *DashboardController {
	return &DashboardController{
		cronService:     cronService,
		executorService: executorService,
	}
}

type StatsResponse struct {
	Tasks     int64 `json:"tasks"`
	Scripts   int64 `json:"scripts"`
	Envs      int64 `json:"envs"`
	Logs      int64 `json:"logs"`
	Scheduled int   `json:"scheduled"`
	Running   int   `json:"running"`
}

func (dc *DashboardController) GetStats(c *gin.Context) {
	var taskCount, scriptCount, envCount, logCount int64

	database.DB.Model(&models.Task{}).Count(&taskCount)
	database.DB.Model(&models.Script{}).Count(&scriptCount)
	database.DB.Model(&models.EnvironmentVariable{}).Count(&envCount)
	database.DB.Model(&models.TaskLog{}).Count(&logCount)

	stats := StatsResponse{
		Tasks:     taskCount,
		Scripts:   scriptCount,
		Envs:      envCount,
		Logs:      logCount,
		Scheduled: dc.cronService.GetScheduledCount(),
		Running:   dc.executorService.GetRunningCount(),
	}

	utils.Success(c, stats)
}

// GetSentence 获取随机古诗词
func (dc *DashboardController) GetSentence(c *gin.Context) {
	utils.Success(c, gin.H{
		"sentence": constant.GetRandomSentence(),
	})
}
