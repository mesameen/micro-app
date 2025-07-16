run:
	cd cmd && go run .

build:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o main cmd/*.go

dockerbuild:
	eval $(minikube docker-env) && docker build -t movie:latest .

dockerpush: 
	minikube image load movie:latest

deploy: dockerpush dockerbuild
	kubectl apply -f kubernetes-deployment.yaml
