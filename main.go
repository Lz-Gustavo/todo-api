package main

import (
	"context"
	"log"
)

var (
	mongoInstance *MongoInstance
)

func init() {
	mongoInstance = &MongoInstance{}
	err := mongoInstance.Connect(context.Background())
	if err != nil {
		log.Fatalln("could not connect to mongo, err:", err.Error())
	}
}

func main() {
	// TODO
}
