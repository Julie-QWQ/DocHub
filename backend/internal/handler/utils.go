package handler

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

const (
	defaultPage     = 1
	defaultPageSize = 20
	maxPageSize     = 100
)

// GetPageParams 从请求中获取分页参数
func GetPageParams(c *gin.Context) (page, pageSize int) {
	page = defaultPage
	pageSize = defaultPageSize

	if p := c.Query("page"); p != "" {
		if parsedPage, err := parseIntParam(p, 1, 100000); err == nil {
			page = parsedPage
		}
	}

	if ps := c.Query("page_size"); ps != "" {
		if parsedSize, err := parseIntParam(ps, 1, maxPageSize); err == nil {
			pageSize = parsedSize
		}
	}

	return page, pageSize
}

func parseIntParam(s string, min, max int) (int, error) {
	var val int
	if _, err := fmt.Sscanf(s, "%d", &val); err != nil {
		return 0, err
	}
	if val < min {
		val = min
	}
	if val > max {
		val = max
	}
	return val, nil
}
