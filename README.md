# react-ts-go-chat-app
High level
- Live chat app with upvoting/downvoting
- Implements sign ups and logins with basic username/password authorization
- Uses an http server for logins and chat history
- Uses a websocket server for live chat and voting updates
- Uses a simple in-memory queue to process live voting updates and avoid race conditions
- React/TS Frontend, Go Backend, Postgres DB running in minikube
# Running locally
Assumes the following are installed
- Docker
- minikube
- npm

First, make sure Docker is running and start minikube in the command line:
```
minikube start
```

### Running the db in minikube
Next, build and start the db:
```
eval $(minikube docker-env) && make k8s-db-rebuild && make k8s-db-run
```

### Running the backend server in minikube
First, from the db minikube output
- Copy the first url (with target port 5432)
- Go to `server.env`
- Update `DB_HOST` formatted as `XXX.XXX.XX.X` from the url
- Update `DB_PORT` to be the port from the url

Next, build and start the server:
```
eval $(minikube docker-env) && make k8s-server-rebuild && make k8s-server-run
```

### Running the frontent client locally
First, from the server minikube output
- Copy the port from the first 127.0.0.1 url and replace the port 4000 in all usages of 127.0.0.1:4000
- Copy the port from the second 127.0.0.1 url and replace the port 4000 in all usages of 127.0.0.1:4001

Next, build and run the client (locally, not minikube)
```
cd client && npm init && npm start
```