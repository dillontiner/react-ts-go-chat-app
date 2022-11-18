# react-ts-go-chat-app

### Running locally with minikube
Assumes the following are installed
- Docker
- minikube
- npm

First, make sure Docker is running and start minikube in the command line:
```
minikube start
```

Next, build and start the db:
```
make k8s-db-rebuild && make k8s-db-run
```
If things fail here, try running `eval $(minikube docker-env)` and running the commands from the make commands individually.

Then
- Copy the first url (with target port 5432)
- Go to `server.env`
- Update `DB_HOST` formatted as `XXX.XXX.XX.X` from the url
- Update `DB_PORT` to be the port from the url

Next, build and start the server:
```
make k8s-server-rebuild && make k8s-server-run
```
If things fail here, try running `eval $(minikube docker-env)` and running the commands from the make commands individually.

Then
- Copy the first url (with http/4000 target port) and replace all usages of 127.0.0.1:4000
- Copy the first url (with websocket/4001 target port) and replace the usage of 127.0.0.1:4001

Next, build and run the client (locally, not minikube)
```
cd client && npm init && npm start
```