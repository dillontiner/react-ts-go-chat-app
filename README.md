# react-ts-go-chat-app

### Running locally with minikube
Assumes the following are installed
- Docker
- minikube

First, build and start the db:
```
make k8s-db-rebuild && make k8s-db-run
```

Then
- Copy the first url (with target port 5432)
- Go to `server.env`
- Update `DB_HOST` formatted as `XXX.XXX.XX.X` from the url
- Update `DB_PORT` to be the port from the url

Next, build and start the server:
```
make k8s-server-rebuild && make k8s-server-run
```