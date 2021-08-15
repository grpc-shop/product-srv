package main

import (
	"flag"
	"fmt"
	"github.com/grpc-shop/product-srv/conf"
	"github.com/grpc-shop/product-srv/handler"
	"github.com/grpc-shop/product-srv/proto/product"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

func main() {
	//读取配置
	dbFile := conf.InitFileConf("../conf/db.json")
	dbConf, err := dbFile.GetDbConf()
	if err != nil {
		fmt.Println("xxx")
		log.Fatal(err)
	}

	db, err := conf.GetDb(dbConf)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("数据库gorm:", db)

	port := flag.Int("port", 0, "the server port")
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatal(err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)

	product.RegisterProductServer(grpcServer, handler.InitProductHandler(db))
	reflection.Register(grpcServer)
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatal(err)
	}
}
