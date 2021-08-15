package model

import (
	"gorm.io/gorm"
)

type ProductState uint32

const (
	ProductInSale ProductState = iota + 1
	ProductNotSale
)

type Product struct {
	Id          int64        `gorm:"column:id" json:"id"`
	Title       string       `gorm:"column:title" json:"title"`
	Description string       `gorm:"column:description" json:"description"`
	Image       string       `gorm:"column:image" json:"image"`
	OnSale      uint32       `gorm:"column:on_sale" json:"onSale"`
	SoldCount   uint32       `gorm:"column:sold_count" json:"soldCount"`
	Price       uint32       `gorm:"column:price" json:"price"`
	CreatedAt   int64        `gorm:"column:created_at;autoCreatedAt" json:"createdAt"`
	UpdatedAt   int64        `gorm:"column:updated_at;autoUpdatedAt" json:"updatedAt"`
	Skus        []ProductSku `gorm:"foreignKey:product_id"`
}

func (Product) TableName() string {
	return "products"
}

func GetWithId(id int64) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id=?", id)
	}
}

func GetWithOnSale(state ProductState) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("on_sale=?", state)
	}
}

func GetWithGreaterCreateTime(time int64) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("created_at>=?", time)
	}
}

func GetWithLessThanCreateTime(time int64) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("created_at<=?", time)
	}
}
