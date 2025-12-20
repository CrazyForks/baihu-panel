package utils

import (
	"strconv"

	"baihu/internal/constant"

	"github.com/gin-gonic/gin"
)

// Pagination 分页参数
type Pagination struct {
	Page     int
	PageSize int
}

// ParsePagination 从请求中解析分页参数
func ParsePagination(c *gin.Context) Pagination {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", strconv.Itoa(constant.DefaultPageSize)))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = constant.DefaultPageSize
	}

	return Pagination{Page: page, PageSize: pageSize}
}

// Offset 计算偏移量
func (p Pagination) Offset() int {
	return (p.Page - 1) * p.PageSize
}

// PaginatedResponse 分页响应
func PaginatedResponse(c *gin.Context, data interface{}, total int64, p Pagination) {
	Success(c, gin.H{
		"data":      data,
		"total":     total,
		"page":      p.Page,
		"page_size": p.PageSize,
	})
}
