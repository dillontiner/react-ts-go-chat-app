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
eval $(minikube docker-env) && make k8s-db-rebuild && make k8s-db-run
```

Then
- Copy the first url (with target port 5432)
- Go to `server.env`
- Update `DB_HOST` formatted as `XXX.XXX.XX.X` from the url
- Update `DB_PORT` to be the port from the url

Next, build and start the server:
```
eval $(minikube docker-env) && make k8s-server-rebuild && make k8s-server-run
```

Then
- Copy the port from the first 127.0.0.1 url and replace the port 4000 in all usages of 127.0.0.1:4000
- Copy the port from the second 127.0.0.1 url and replace the port 4000 in all usages of 127.0.0.1:4001

Next, build and run the client (locally, not minikube)
```
cd client && npm init && npm start
```