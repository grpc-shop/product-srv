package dao

import (
	"github.com/grpc-shop/product-srv/model"
	"github.com/grpc-shop/product-srv/param"
	"gorm.io/gorm"
)

type ProductDao interface {
	GetProductList(param param.GetListParam) (list []*model.Product, count int64, err error)
	CreateProduct(product *model.Product) (id int64, err error)
	UpdateProductById(product *model.Product, id int64) error
	FirstProductById(id int64) (*model.Product, error)
	DeleteProductById(id int64) (count int64, err error)
	ListingProductById(id int64) (count int64, err error)
	DeListingProductById(id int64) (count int64, err error)
}

var _ ProductDao = (*ProductImpl)(nil)

type ProductImpl struct {
	db *gorm.DB
}

func NewProductImpl(db *gorm.DB) ProductDao {
	return &ProductImpl{
		db: db,
	}
}

func (p *ProductImpl) GetProductList(param param.GetListParam) (list []*model.Product, count int64, err error) {
	base := p.db.Model(&model.Product{}).Debug().
		Preload("Skus", func(db *gorm.DB) *gorm.DB {
			return db.Select(
				"id,title,price,stock,product_id,created_at")
		})

	var (
		scopes []func(db *gorm.DB) *gorm.DB
	)

	scopes = append(scopes, model.GetWithOnSale(model.ProductInSale))
	if param.StartCreateTime > 0 {
		scopes = append(scopes, model.GetWithGreaterCreateTime(param.StartCreateTime))
	}

	if param.EndCreateTime > 0 {
		scopes = append(scopes, model.GetWithLessThanCreateTime(param.EndCreateTime))
	}

	base.Scopes(scopes...)

	base.Count(&count)

	err = base.Scopes(param.Page.Paginate()).Find(&list).Error
	return
}

func (p *ProductImpl) CreateProduct(product *model.Product) (id int64, err error) {
	err = p.db.Model(&model.Product{}).Create(product).Error
	id = product.Id
	return
}

func (p *ProductImpl) UpdateProductById(product *model.Product, id int64) error {
	return p.db.Model(&model.Product{}).
		Scopes(model.GetWithId(id)).
		Updates(product).Error
}

func (p *ProductImpl) DeleteProductById(id int64) (count int64, err error) {
	base := p.db.Model(&model.Product{}).
		Select("Skus").
		Delete(&model.Product{}, id)
	count = base.RowsAffected
	err = base.Error
	return
}

func (p *ProductImpl) FirstProductById(id int64) (*model.Product, error) {
	var (
		product model.Product
	)
	err := p.db.Model(model.Product{}).First(&product, id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &product, nil
}

func (p *ProductImpl) ListingProductById(id int64) (count int64, err error) {
	base := p.db.Model(&model.Product{}).
		Where("id=?", id).
		Update("on_sale", model.ProductInSale)
	count = base.RowsAffected
	err = base.Error
	return
}

func (p *ProductImpl) DeListingProductById(id int64) (count int64, err error) {
	base := p.db.Model(&model.Product{}).
		Where("id=?", id).
		Update("on_sale", model.ProductNotSale)
	count = base.RowsAffected
	err = base.Error
	return
}
