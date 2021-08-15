server:
	cd cmd && go run main.go -port 8899
proto:
	protoc -I=proto proto/*.proto --go_out=:proto --go-grpc_out=:proto \
-I=${GOPATH}/pkg/mod/github.com/envoyproxy/protoc-gen-validate@v0.6.1 \
--validate_out="lang=go:./proto"\

lint:
	golangci-lint run

.PHONY:server proto



###########git############

pull:
	git push --rebase

checkout:
	@# echo $x;
	- git checkout $x;





## https://studygolang.com/articles/25743
## https://www.cnblogs.com/yisany/p/14888041.html