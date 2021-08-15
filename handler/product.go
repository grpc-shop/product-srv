package handler

import (
	"context"
	"github.com/grpc-shop/product-srv/model"
	"github.com/grpc-shop/product-srv/param"
	"github.com/grpc-shop/product-srv/proto/product"
	"github.com/grpc-shop/product-srv/service"
)

var _ product.ProductServer = (*ProductHandler)(nil)

type ProductHandler struct {
	product.UnimplementedProductServer
	server service.ProductServer
}

func NewProductHandler(server service.ProductServer) *ProductHandler {
	return &ProductHandler{server: server}
}

func (p *ProductHandler) GetProductList(ctx context.Context, req *product.GetProductListReq) (*product.GetProductListReply, error) {
	resp := product.GetProductListReply{
		Code: product.Code_Success,
		Data: new(product.GetProductListReplyProduct),
	}
	getListParam := param.InitGetListParam(
		req.GetStartCreateTime(),
		req.GetEndCreateTime(), req.GetPage(), req.GetPageSize())

	list, count, err := p.server.GetProductList(getListParam)
	if err != nil {
		resp.Code = product.Code_GetProductErr
		return nil, err
	}
	for index := range list {
		productModel := list[index]
		productPb := product.ProductEntity{
			Id:          productModel.Id,
			Title:       productModel.Title,
			Description: productModel.Description,
			Image:       productModel.Image,
			OnSale:      productModel.OnSale,
			SoldCount:   productModel.SoldCount,
			Price:       productModel.Price,
			CreatedAt:   productModel.CreatedAt,
		}

		skus := productModel.Skus
		for j := range skus {
			skuPb := product.ProductSku{
				Id:        skus[j].Id,
				Title:     skus[j].Title,
				Price:     skus[j].Price,
				Stock:     skus[j].Stock,
				ProductId: skus[j].ProductId,
				CreatedAt: skus[j].CreatedAt,
			}
			productPb.Skus = append(productPb.Skus, &skuPb)
		}

		resp.Data.ProductList = append(resp.Data.ProductList, &productPb)
	}

	pageInfo := getListParam.Page.GetListPage(count)

	resp.Data.Total = pageInfo.Total
	resp.Data.PageSize = pageInfo.PageSize
	resp.Data.Page = pageInfo.Page
	resp.Data.TotalPage = pageInfo.TotalPage
	return &resp, nil
}

func (p *ProductHandler) CreateProduct(ctx context.Context, req *product.CreateProductReq) (*product.CreateProductReply, error) {
	resp := product.CreateProductReply{
		Code: product.Code_Success,
		Data: new(product.CreateProductReplyProduct),
	}
	err := req.Validate()
	if err != nil {
		return &resp, err
	}

	var (
		productModel model.Product
	)
	productModel.Title = req.GetProduct().GetTitle()
	productModel.Description = req.GetProduct().GetDescription()
	productModel.Image = req.GetProduct().GetImage()
	productModel.OnSale = req.GetProduct().GetOnSale()

	for i := range req.Product.Skus {
		var (
			skuModel model.ProductSku
		)
		skuReq := req.Product.GetSkus()[i]
		skuModel.Price = skuReq.GetPrice()
		skuModel.Title = skuReq.GetTitle()
		skuModel.Stock = skuReq.GetStock()
		productModel.Skus = append(productModel.Skus, skuModel)
	}

	min := productModel.Skus[0].Price

	for i := range productModel.Skus {
		if productModel.Skus[i].Price < min {
			min = productModel.Skus[i].Price
		}
	}
	productModel.Price = min

	id, err := p.server.CreateProduct(&productModel)
	if err != nil {
		return nil, err
	}
	resp.Data.Id = id
	return &resp, nil
}

func (p *ProductHandler) UpdateProduct(ctx context.Context, req *product.UpdateProductReq) (*product.UpdateProductReply, error) {
	resp := product.UpdateProductReply{
		Code: product.Code_Success,
	}

	err := req.Validate()
	if err != nil {
		return &resp, err
	}

	var (
		productModel model.Product
	)
	productModel.Title = req.GetProduct().GetTitle()
	productModel.Description = req.GetProduct().GetDescription()
	productModel.Image = req.GetProduct().GetImage()
	productModel.OnSale = req.GetProduct().GetOnSale()

	for i := range req.Product.Skus {
		var (
			skuModel model.ProductSku
		)
		skuReq := req.Product.GetSkus()[i]
		skuModel.Price = skuReq.GetPrice()
		skuModel.Title = skuReq.GetTitle()
		skuModel.Stock = skuReq.GetStock()
		productModel.Skus = append(productModel.Skus, skuModel)
	}
	min := productModel.Skus[0].Price

	for i := range productModel.Skus {
		if productModel.Skus[i].Price < min {
			min = productModel.Skus[i].Price
		}
	}
	productModel.Price = min

	err = p.server.UpdateProduct(req.GetId(), &productModel)
	if err != nil {
		resp.Code = product.Code_UpdateProductErr
		return &resp, err
	}
	return &resp, nil
}

func (p *ProductHandler) DeleteProduct(ctx context.Context, req *product.DeleteProductReq) (*product.DeleteProductReply, error) {
	resp := product.DeleteProductReply{
		Code: product.Code_Success,
		Data: new(product.DeleteProductReplyProduct),
	}
	err := req.Validate()
	if err != nil {
		return &resp, err
	}

	ids, err := p.server.DeleteProductByIds(req.GetIds())
	if err != nil {
		resp.Code = product.Code_DeleteProductErr
		return &resp, err
	}
	resp.Data.Ids = ids
	return &resp, nil
}

func (p *ProductHandler) ListingProduct(ctx context.Context, req *product.ListingProductReq) (*product.ListingProductReply, error) {
	resp := product.ListingProductReply{
		Code: product.Code_Success,
		Data: new(product.ListingProductReplyProduct),
	}
	err := req.Validate()
	if err != nil {
		return &resp, err
	}

	ids, err := p.server.ListingProductById(req.GetIds())
	if err != nil {
		resp.Code = product.Code_ListingProductErr
		return &resp, err
	}
	resp.Data.Ids = ids
	return &resp, nil
}

func (p *ProductHandler) DeListingProduct(ctx context.Context, req *product.DeListingProductReq) (*product.DeListingProductReply, error) {
	resp := product.DeListingProductReply{
		Code: product.Code_Success,
		Data: new(product.DeListingProductReplyProduct),
	}
	err := req.Validate()
	if err != nil {
		return &resp, err
	}

	ids, err := p.server.DeListingProductById(req.GetIds())
	if err != nil {
		resp.Code = product.Code_DeListingProductErr
		return &resp, err
	}
	resp.Data.Ids = ids
	return &resp, nil
}
