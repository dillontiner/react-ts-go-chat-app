docker-db:
	docker run \
    -it \
    --rm \
    -p 5432:5432 \
    chat-room-db