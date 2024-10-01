.PHONY: generate_go

generate_go:
	protoc --go_out=pb/gen --go-grpc_out=pb/gen ./pb/payment.proto

# naming to gcr.io/$(PROJECT_ID)/grpc-payment-service:v1
docker_build_grpc:
	docker build -t gcr.io/$(PROJECT_ID)/grpc-payment-service:v1 -f build/Dockerfile.grpc .

# naming to gcr.io/$(PROJECT_ID)/rest-payment-service:v1  
docker_build_rest:
	docker build -t gcr.io/$(PROJECT_ID)/rest-payment-service:v1 -f build/Dockerfile.rest .

docker_push_grpc:
	docker push gcr.io/$(PROJECT_ID)/grpc-payment-service:v1

docker_push_rest:
	docker push gcr.io/$(PROJECT_ID)/rest-payment-service:v1

gcloud_cluster_create:
	gcloud container clusters create payment-cluster --num-nodes=3 --location=europe-west1-b

restart_cluster:
	kubectl delete deployment grpc-payment-service
	kubectl delete deployment rest-payment-service
	kubectl delete service grpc-payment-service
	kubectl delete service rest-payment-service
	kubectl apply -f grpc-deployment.yaml
	kubectl apply -f grpc-service.yaml
	kubectl apply -f rest-deployment.yaml
	kubectl apply -f rest-service.yaml
	kubectl get pods

get_services:
	kubectl get services

get_deployments:
	kubectl describe pod rest-payment-service-

get_pods:
	kubectl get pods -o wide

get_nodes:
	kubectl get nodes -o wide

connect_to_vm:
	gcloud compute ssh grpc-vm --zone=europe-west4-a