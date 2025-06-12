package helper

import (
	"github.com/gin-gonic/gin"
	"github.com/gogf/gf/v2/util/gconv"
	"gorm.io/gorm"
)

type Pagination struct {
	Page int
	PageSize int
	Total int64
	WithoutCount bool
}
func (p *Pagination) Build(db *gorm.DB, c *gin.Context) *gorm.DB {
	if c != nil {
		p.Page = gconv.Int(c.DefaultQuery("pagination[page]", "1"))	
		p.PageSize = gconv.Int(c.DefaultQuery("pagination[pageSize]", "20"))	
		p.WithoutCount = gconv.Bool(c.DefaultQuery("pagination[withoutCount]", "0"))
	}

	if p.Page <= 0 {
		p.Page = 1
	}
	if p.PageSize <= 0 {
		p.PageSize = 30
	}
	// 查询总数
	if !p.WithoutCount {
		db.Count(&p.Total)
	}
	return db.Limit(p.PageSize).Offset(p.PageSize * (p.Page - 1))
}