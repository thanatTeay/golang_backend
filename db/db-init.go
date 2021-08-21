package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Resource struct {
	DB *mongo.Database
}

func (r *Resource) Close() {
	logrus.Warning("Closing all db connections")
}

func InitResourc() (*Resource, error) {
	dbName := "UserData"
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://Thanat2208:thegame901@crud.qooys.mongodb.net/UserData?retryWrites=true&w=majority"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	err = client.Connect(ctx)
	//log.Fatal(cancel)
	if err != nil {
		log.Fatal(err)
	}
	defer cancel()
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	databases, err := client.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(databases)
	//log.Fatal(client.Database("UserData"))

	return &Resource{DB: client.Database(dbName)}, nil

}
