package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"github.com/grpc-shop/product-srv/conf"
	"github.com/grpc-shop/product-srv/handler"
	"github.com/grpc-shop/product-srv/proto/product"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
	"io/ioutil"
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
	port := flag.Int("port", 0, "the server port")
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatal(err)
	}

	// ctl
	certificate, err := tls.LoadX509KeyPair("../tool/cert/server-cert.pem", "../tool/cert/server-key.pem")
	if err != nil {
		log.Fatal(err)
	}

	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile("../tool/cert/ca-cert.pem")
	if err != nil {
		log.Fatal(err)
	}
	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		log.Fatal("failed to append certs")
	}
	creds := credentials.NewTLS(&tls.Config{
		Certificates:       []tls.Certificate{certificate},
		ClientAuth:         tls.RequireAndVerifyClientCert, // NOTE: this is optional!
		ClientCAs:          certPool,
	})
	var opts []grpc.ServerOption

	opts = append(opts, grpc.Creds(creds))

	grpcServer := grpc.NewServer(opts...)

	product.RegisterProductServer(grpcServer, handler.InitProductHandler(db))
	reflection.Register(grpcServer)

	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatal(err)
	}
}
