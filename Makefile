docker-db-build:
	docker build -t chat-app-db ./db

docker-db-run:
	docker run \
		-it \
		--rm \
		-p 5432:5432 \
		chat-app-db

k8s-db-delete:
	kubectl delete deployment chat-app-db && kubectl delete service chat-app-db && kubectl delete pod chat-app-db

k8s-db-build:
	kubectl apply -f ./db/db-k8s-deployment.yml && \
	kubectl apply -f ./db/db-k8s-service.yml

k8s-db-rebuild:
	eval $(minikube docker-env) && \
	make k8s-db-delete || true && \
	make docker-db-build && \
	make k8s-db-build

k8s-db-run:
	minikube service chat-app-db

docker-server-build:
	docker build -t chat-app-server ./server

docker-server-run:
	docker run \
		-it \
		--rm \
		-v ${PWD}/server:/app \
		-p 4000:4000 \
		-p 4001:4001 \
		--network="host" \
		chat-app-server

k8s-server-delete:
	kubectl delete deployment chat-app-server && kubectl delete service chat-app-server && kubectl delete pod chat-app-server

k8s-server-build:
	kubectl apply -f ./server/server-k8s-deployment.yml && \
	kubectl apply -f ./server/server-k8s-service.yml

k8s-server-rebuild:
	eval $(minikube docker-env) && \
	make k8s-server-delete || true && \
	make docker-server-build && \
	make k8s-server-build

k8s-server-run:
	minikube service chat-app-server