package service

import (
	"errors"
	"fmt"
	"github.com/wuqinqiang/grpc-shop/product-srv/dao"
	"github.com/wuqinqiang/grpc-shop/product-srv/model"
	"github.com/wuqinqiang/grpc-shop/product-srv/param"
)

var (
	ProductNoFoundErr = errors.New("订单不存在")
)

type ProductServer interface {
	GetProductList(param param.GetListParam) (list []*model.Product, count int64, err error)
	CreateProduct(product *model.Product) (id int64, err error)
	UpdateProduct(id int64, product *model.Product) error
	DeleteProductByIds(ids []int64) ([]int64, error)
	ListingProductById(ids []int64) ([]int64, error)
	DeListingProductById(ids []int64) ([]int64, error)
}

var _ ProductServer = (*ProductServerImpl)(nil)

type ProductServerImpl struct {
	productDao dao.ProductDao
}

func NewProductServerImpl(dao dao.ProductDao) ProductServer {
	return &ProductServerImpl{productDao: dao}
}

func (s *ProductServerImpl) GetProductList(param param.GetListParam) (list []*model.Product, count int64, err error) {
	return s.productDao.GetProductList(param)
}

func (s *ProductServerImpl) CreateProduct(product *model.Product) (id int64, err error) {
	return s.productDao.CreateProduct(product)
}

func (s *ProductServerImpl) UpdateProduct(id int64, product *model.Product) error {
	firstProduct, err := s.productDao.FirstProductById(id)
	if err != nil {
		return err
	}
	if firstProduct == nil {
		return ProductNoFoundErr
	}
	return s.productDao.UpdateProductById(product, id)
}

func (s *ProductServerImpl) DeleteProductByIds(ids []int64) ([]int64, error) {
	var (
		successIds []int64
	)
	for _, id := range ids {
		firstProduct, err := s.productDao.FirstProductById(id)
		if err != nil || firstProduct == nil {
			fmt.Println("商品id:不存在，不能删除", id)
			continue
		}
		if model.ProductState(firstProduct.OnSale) == model.ProductInSale {
			fmt.Println("商品id:处于上架状态，不能删除", id)
			continue
		}
		count, err := s.productDao.DeleteProductById(id)
		if err != nil {
			fmt.Printf("商品id:%v删除失败:%v", id, err)
		}
		if count > 0 {
			successIds = append(successIds, id)
		}
	}
	return successIds, nil
}

func (s *ProductServerImpl) ListingProductById(ids []int64) ([]int64, error) {
	var (
		successIds []int64
	)
	for _, id := range ids {
		firstProduct, err := s.productDao.FirstProductById(id)
		if err != nil || firstProduct == nil {
			fmt.Println("商品id:不存在", id)
			continue
		}
		if model.ProductState(firstProduct.OnSale) == model.ProductInSale {
			fmt.Println("商品id:处于上架状态，请勿重复操作", id)
			continue
		}
		count, err := s.productDao.ListingProductById(id)
		if err != nil {
			fmt.Printf("商品id:%v上架失败:%v", id, err)
		}
		if count > 0 {
			successIds = append(successIds, id)
		}
	}
	return successIds, nil
}

func (s *ProductServerImpl) DeListingProductById(ids []int64) ([]int64, error) {
	var (
		successIds []int64
	)
	for _, id := range ids {
		firstProduct, err := s.productDao.FirstProductById(id)
		if err != nil || firstProduct == nil {
			fmt.Println("商品id:不存在", id)
			continue
		}
		if model.ProductState(firstProduct.OnSale) == model.ProductNotSale {
			fmt.Println("商品id:处于下架状态，请勿重复操作", id)
			continue
		}
		count, err := s.productDao.DeListingProductById(id)
		if err != nil {
			fmt.Printf("商品id:%v下架失败:%v", id, err)
		}
		if count > 0 {
			successIds = append(successIds, id)
		}
	}
	return successIds, nil
}
