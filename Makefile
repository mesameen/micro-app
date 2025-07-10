proto:
	cd src/api && protoc *.proto \
		--go_out=. \
		--go-grpc_out=. \
		--proto_path=.
