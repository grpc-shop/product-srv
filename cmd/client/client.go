package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"github.com/grpc-shop/product-srv/proto/product"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	"log"
)

func main() {
	serverAddress := flag.String("address", "", "the server address")
	flag.Parse()
	log.Printf("dial server %s", *serverAddress)

	//ctl
	certificate, err := tls.LoadX509KeyPair("../../tool/cert/client-cert.pem", "../../tool/cert/client-key.pem")
	if err != nil {
		log.Fatal("init certificate err:", err)
	}
	certPool := x509.NewCertPool()

	ca, err := ioutil.ReadFile("../../tool/cert/ca-cert.pem")
	if err != nil {
		log.Fatal("read ca.crt err:", err)
	}
	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		log.Fatal("failed to append ca certs")
	}

	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{certificate},
		RootCAs:      certPool,
	})
	cc, err := grpc.Dial(*serverAddress, grpc.WithTransportCredentials(creds))
	//cc, err := grpc.Dial(*serverAddress, grpc.WithInsecure())
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
