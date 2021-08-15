package model

type ProductSku struct {
	Id          uint64 `gorm:"column:id" json:"id"`
	Title       string `gorm:"column:title" json:"title"`
	Price       uint32 `gorm:"column:price" json:"price"`
	Stock       uint32 `gorm:"column:stock" json:"stock"`
	ProductId   uint64 `gorm:"column:product_id" json:"productId"`
	CreatedAt   int64  `gorm:"column:created_at;autoCreatedAt" json:"createdAt"`
	UpdatedAt   int64  `gorm:"column:updated_at;autoUpdatedAt" json:"updatedAt"`
}

func (ProductSku) TableName() string {
	return "product_skus"
}
