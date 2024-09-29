.PHONY: generate_go

generate_go:
	protoc --go_out=pb/gen --go-grpc_out=pb/gen ./pb/payment.proto

# naming to gcr.io/916153668772/grpc-payment-service:v1
docker_build_grpc:
	docker build -t gcr.io/916153668772/grpc-payment-service:v1 -f build/Dockerfile.grpc .

# naming to gcr.io/916153668772/rest-payment-service:v1 
docker_build_rest:
	docker build -t gcr.io/916153668772/rest-payment-service:v1 -f build/Dockerfile.rest .