package db

import (
	"gorm.io/gorm"
	"math"
)

const (
	DefaultSize    = 20
	DefaultPage    = 1
	DefaultMaxSize = 50
)

type Page struct {
	Page     int64 `json:"page"`
	PageSize int64 `json:"pageSize"`
}

func InitPageSize(page, size int64) Page {
	if page <= 0 {
		page = DefaultPage
	}
	if size <= 0 {
		size = DefaultSize
	}
	if size > DefaultMaxSize {
		size = DefaultMaxSize
	}
	return Page{page, size}
}

func (p *Page) Offset() int64 {
	return p.PageSize * (p.Page - 1)
}

func (p *Page) GetListPage(total int64) ListPage {
	var (
		res ListPage
	)
	res.TotalPage = int64(math.Ceil(float64(total) / float64(p.PageSize)))
	res.PageSize = p.PageSize
	res.Total = total
	res.Page = p.Page
	return res
}

func (p *Page) Paginate() func(d *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Limit(int(p.PageSize)).Offset(int(p.Offset()))
	}
}

type ListPage struct {
	TotalPage int64 `json:"totalPage"`
	Page      int64 `json:"page"`
	PageSize  int64 `json:"pageSize"`
	Total     int64 `json:"total"`
}
