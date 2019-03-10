package app

import (
	"context"
	"github.com/A9u/function_junction/config"
	"github.com/mongodb/mongo-go-driver/mongo"
	"go.uber.org/zap"
	"time"
)

var (
	db     *mongo.Database
	client *mongo.Client
	logger *zap.SugaredLogger
	ctx    context.Context
)

func Init() {
	InitLogger()

	err := initDB()
	if err != nil {
		panic(err)
	}
}

func InitLogger() {
	zapLogger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	logger = zapLogger.Sugar()
}

func initDB() (err error) {
	dbConfig := config.Database()
	client, err := mongo.NewClient(dbConfig.ConnectionURL())
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		return
	}
	db = client.Database(dbConfig.DbName())
	/*
		if err = client.Ping(ctx, nil); err != nil {
			return
		}
	*/

	return
}

func GetCollection(name string) *mongo.Collection {
	collection := db.Collection(name)
	return collection
}

func GetDB() *mongo.Database {
	return db
}

func GetLogger() *zap.SugaredLogger {
	return logger
}

func Close() {
	logger.Sync()
	client.Disconnect(nil)
}
