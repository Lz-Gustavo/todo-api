FROM golang:1.15-alpine as build-env
ENV GO111MODULE=on

RUN mkdir /go/src/todo-api
WORKDIR /go/src/todo-api

COPY go.mod .
COPY go.sum .
#RUN apk update && apk add git && apk add gcc

RUN go mod download
COPY . .
RUN go build -o todoapi

FROM alpine:latest
WORKDIR /root/
COPY --from=build-env /go/src/todo-api/todoapi .
ENV GIPHY_APIKEY=aXQy4rfovFa6J18bXZHt6MKJ8hQQjkXd

CMD ["./todoapi"]