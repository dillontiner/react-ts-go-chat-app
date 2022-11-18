FROM golang:1.18.4-alpine

WORKDIR /server

ADD ./ ./

RUN go mod download

RUN go mod tidy

RUN go build -o /chat-app-server

EXPOSE 4000
EXPOSE 4001

CMD ["/chat-app-server"]