proto:
	cd src/api && protoc *.proto \
		--go_out=. \
		--go-grpc_out=. \
		--proto_path=.

build:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o main rating/cmd/*.go

run:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go run rating/cmd/*.go

dockerbuild:
	eval $(minikube docker-env) && docker build -t rating:latest .
