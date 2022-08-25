package dal

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"kp-management/internal/pkg/conf"
)

var (
	m *mongo.Client
)

func MustInitMongo() {
	var err error
	clientOptions := options.Client().ApplyURI(conf.Conf.MongoDB.DSN)

	m, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		panic(fmt.Errorf("mongo err:%w", err))
	}

	err = m.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("mongodb initialized")
}

func GetMongo() *mongo.Client {
	return m
}

func MongoDB() string {
	return conf.Conf.MongoDB.Database
}
