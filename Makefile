proto:
	protoc -I. pkg/protos/**/*.proto --go_out=:. --go-grpc_out=:. --go-grpc_opt=module=github.com/sumit-tembe/grpc-svc --go_opt=module=github.com/sumit-tembe/grpc-svc

