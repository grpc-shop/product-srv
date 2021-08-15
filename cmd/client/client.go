package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/wuqinqiang/grpc-shop/product-srv/proto/product"
	"google.golang.org/grpc"
	"log"
)

func main() {
	serverAddress := flag.String("address", "", "the server address")
	flag.Parse()
	log.Printf("dial server %s", *serverAddress)

	cc, err := grpc.Dial(*serverAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatal("cannot dial server:", err)
	}

	lpClient := product.NewProductClient(cc)
	list, err := lpClient.GetProductList(context.Background(), &product.GetProductListReq{})
	if err != nil {
		fmt.Println("错误:", err)
		return
	}
	fmt.Println("数据:", list)
}
