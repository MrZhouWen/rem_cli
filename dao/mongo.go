package dao

import (
	"context"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

var Mongo *mongo.Client
var MongoDb *mongo.Database

func InitMongo() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	o := options.Client().ApplyURI(viper.GetString("mongodb.uri"))
	o.SetMaxPoolSize(20)
	o.SetReplicaSet(viper.GetString("mongodb.replicaSet"))
	o.SetReadPreference(readpref.PrimaryPreferred())
	var err error
	Mongo, err = mongo.Connect(ctx, o)
	if err != nil {
		panic(err)
	}
	MongoDb = Mongo.Database("rem")

	err = Mongo.Ping(ctx, nil)
	if err != nil {
		panic(err)
	}
}
