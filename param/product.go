package param

import "github.com/grpc-shop/product-srv/tool/db"

type GetListParam struct {
	StartCreateTime int64
	EndCreateTime   int64
	Page            db.Page
}

func InitGetListParam(start, end int64, page, pageSize int64) GetListParam {
	return GetListParam{
		StartCreateTime: start,
		EndCreateTime:   end,
		Page:            db.InitPageSize(page, pageSize),
	}
}
