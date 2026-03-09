package controllers

import (
	"github.com/engigu/baihu-panel/internal/services"
	"github.com/engigu/baihu-panel/internal/utils"

	"github.com/gin-gonic/gin"
)

type AppLogController struct {
	appLogService *services.AppLogService
}

func NewAppLogController() *AppLogController {
	return &AppLogController{
		appLogService: services.NewAppLogService(),
	}
}

// GetLogs 获取应用日志列表
// @Summary 获取应用日志列表
// @Tags 应用日志
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param category query string false "日志分类"
// @Param status query string false "状态"
// @Param level query string false "级别"
// @Param keyword query string false "搜索关键词"
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} utils.Response
// @Router /app-logs [get]
func (ac *AppLogController) GetLogs(c *gin.Context) {
	p := utils.ParsePagination(c)
	category := c.Query("category")
	status := c.Query("status")
	level := c.Query("level")
	keyword := c.Query("keyword")

	logs, total, err := ac.appLogService.List(category, status, level, p.Page, p.PageSize, keyword)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	utils.PaginatedResponse(c, logs, total, p)
}

// MarkAsRead 标记已读
// @Summary 标记已读
// @Tags 应用日志
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body object true "请求体"
// @Success 200 {object} utils.Response
// @Router /app-logs/read [post]
func (ac *AppLogController) MarkAsRead(c *gin.Context) {
	var req struct {
		ID       string `json:"id"`
		Category string `json:"category"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	if req.ID != "" {
		if err := ac.appLogService.MarkAsRead(req.ID); err != nil {
			utils.BadRequest(c, err.Error())
			return
		}
	} else if req.Category != "" {
		if err := ac.appLogService.MarkAllAsRead(req.Category); err != nil {
			utils.BadRequest(c, err.Error())
			return
		}
	} else {
		utils.BadRequest(c, "id 或 category 必须提供")
		return
	}
	utils.SuccessMsg(c, "标记成功")
}

// ClearLogs 清理日志
// @Summary 清理日志
// @Tags 应用日志
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body object true "请求体"
// @Success 200 {object} utils.Response
// @Router /app-logs/clear [post]
func (ac *AppLogController) ClearLogs(c *gin.Context) {
	var req struct {
		Category string `json:"category"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	if err := ac.appLogService.Clear(req.Category); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	utils.SuccessMsg(c, "清理成功")
}
