package controllers

import (
	"strconv"

	"baihu/internal/database"
	"baihu/internal/models"
	"baihu/internal/utils"

	"github.com/gin-gonic/gin"
)

type LogController struct{}

func NewLogController() *LogController {
	return &LogController{}
}

type TaskLogResponse struct {
	ID        uint             `json:"id"`
	TaskID    uint             `json:"task_id"`
	TaskType  string           `json:"task_type"`
	TaskName  string           `json:"task_name"`
	Command   string           `json:"command"`
	Status    string           `json:"status"`
	Duration  int64            `json:"duration"`
	CreatedAt models.LocalTime `json:"created_at"`
}

func (lc *LogController) GetLogs(c *gin.Context) {
	p := utils.ParsePagination(c)
	taskID, _ := strconv.Atoi(c.DefaultQuery("task_id", "0"))
	taskName := c.DefaultQuery("task_name", "")
	taskType := c.DefaultQuery("task_type", "") // task, sync

	var logs []models.TaskLog
	var total int64

	query := database.DB.Model(&models.TaskLog{})
	if taskID > 0 {
		query = query.Where("task_id = ?", taskID)
	}

	// 按任务类型过滤
	if taskType != "" {
		query = query.Where("task_type = ?", taskType)
	}

	// 按任务名称过滤
	if taskName != "" {
		if taskType == "sync" {
			// 同步任务
			var taskIDs []uint
			database.DB.Model(&models.SyncTask{}).Where("name LIKE ?", "%"+taskName+"%").Pluck("id", &taskIDs)
			if len(taskIDs) > 0 {
				query = query.Where("task_id IN ?", taskIDs)
			} else {
				utils.PaginatedResponse(c, []TaskLogResponse{}, 0, p)
				return
			}
		} else if taskType == "task" || taskType == "" {
			// 普通任务或全部
			var taskIDs []uint
			database.DB.Model(&models.Task{}).Where("name LIKE ?", "%"+taskName+"%").Pluck("id", &taskIDs)
			if taskType == "" {
				// 全部类型时，也搜索同步任务
				var syncTaskIDs []uint
				database.DB.Model(&models.SyncTask{}).Where("name LIKE ?", "%"+taskName+"%").Pluck("id", &syncTaskIDs)
				if len(taskIDs) > 0 || len(syncTaskIDs) > 0 {
					query = query.Where("(task_type = 'task' AND task_id IN ?) OR (task_type = 'sync' AND task_id IN ?)", taskIDs, syncTaskIDs)
				} else {
					utils.PaginatedResponse(c, []TaskLogResponse{}, 0, p)
					return
				}
			} else if len(taskIDs) > 0 {
				query = query.Where("task_id IN ?", taskIDs)
			} else {
				utils.PaginatedResponse(c, []TaskLogResponse{}, 0, p)
				return
			}
		}
	}

	query.Count(&total)
	query.Order("id DESC").Offset(p.Offset()).Limit(p.PageSize).Find(&logs)

	// 收集任务ID
	taskIDList := make([]uint, 0)
	syncTaskIDList := make([]uint, 0)
	for _, log := range logs {
		if log.TaskType == "sync" {
			syncTaskIDList = append(syncTaskIDList, log.TaskID)
		} else {
			taskIDList = append(taskIDList, log.TaskID)
		}
	}

	// 获取普通任务名称
	taskMap := make(map[uint]string)
	if len(taskIDList) > 0 {
		var tasks []models.Task
		database.DB.Where("id IN ?", taskIDList).Find(&tasks)
		for _, t := range tasks {
			taskMap[t.ID] = t.Name
		}
	}

	// 获取同步任务名称
	syncTaskMap := make(map[uint]string)
	if len(syncTaskIDList) > 0 {
		var syncTasks []models.SyncTask
		database.DB.Where("id IN ?", syncTaskIDList).Find(&syncTasks)
		for _, t := range syncTasks {
			syncTaskMap[t.ID] = t.Name
		}
	}

	result := make([]TaskLogResponse, len(logs))
	for i, log := range logs {
		taskName := ""
		if log.TaskType == "sync" {
			taskName = syncTaskMap[log.TaskID]
		} else {
			taskName = taskMap[log.TaskID]
		}
		result[i] = TaskLogResponse{
			ID:        log.ID,
			TaskID:    log.TaskID,
			TaskType:  log.TaskType,
			TaskName:  taskName,
			Command:   log.Command,
			Status:    log.Status,
			Duration:  log.Duration,
			CreatedAt: log.CreatedAt,
		}
	}

	utils.PaginatedResponse(c, result, total, p)
}

func (lc *LogController) GetLogDetail(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "无效的日志ID")
		return
	}

	var log models.TaskLog
	if err := database.DB.First(&log, id).Error; err != nil {
		utils.NotFound(c, "日志不存在")
		return
	}

	utils.Success(c, gin.H{
		"id":         log.ID,
		"task_id":    log.TaskID,
		"command":    log.Command,
		"output":     log.Output,
		"status":     log.Status,
		"duration":   log.Duration,
		"created_at": log.CreatedAt,
	})
}
