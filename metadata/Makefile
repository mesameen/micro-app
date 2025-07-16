run:
	cd cmd && go run .

build:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o main cmd/*.go

dockerbuild: build
	eval $(minikube docker-env) && docker build -t metadata:latest .

dockerpush: 
	minikube image load metadata:latest

deploy: dockerbuild dockerpush
	kubectl apply -f kubernetes-deployment.yaml
