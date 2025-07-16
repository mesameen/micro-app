run:
	cd cmd && go run .

build:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -o main cmd/*.go

dockerbuild: build
	eval $(minikube docker-env) && docker build -t rating:latest .

dockerpush: 
	minikube image load rating:latest

deploy: dockerbuild dockerpush
	kubectl apply -f kubernetes-deployment.yaml
