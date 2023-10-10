genIDL:
	protoc --go_out=./pb --go-grpc_out=./pb ./pb/*.proto
build:
	docker build -t tinycache ./