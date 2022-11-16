docker-db-build:
	docker build -t chat-room-db ./db

docker-db-run:
	docker run \
		-it \
		--rm \
		-p 5432:5432 \
		chat-room-db

docker-server-build:
	docker build -t chat-room-server ./server

docker-server-run:
	docker run \
		-it \
		--rm \
		-v ${PWD}/server:/app \
		-p 5000:5000 \
		--network="host" \
		chat-room-server

