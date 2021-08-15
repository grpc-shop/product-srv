// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//+build !wireinject

package handler

import (
	"github.com/grpc-shop/product-srv/dao"
	"github.com/grpc-shop/product-srv/service"
	"gorm.io/gorm"
)

// Injectors from wire.go:

func InitProductHandler(db *gorm.DB) *ProductHandler {
	productDao := dao.NewProductImpl(db)
	productServer := service.NewProductServerImpl(productDao)
	productHandler := NewProductHandler(productServer)
	return productHandler
}
