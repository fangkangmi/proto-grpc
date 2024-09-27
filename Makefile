.PHONY: generate_go

generate_go:
	protoc --go_out=pb/gen --go-grpc_out=pb/gen ./pb/payment.proto