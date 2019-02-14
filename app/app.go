package app

import (
	"fmt"
	// "time"
	"github.com/globalsign/mgo"
	// "github.com/jmoiron/sqlx"
	"github.com/joshsoftware/golang-boilerplate/config"
	"go.uber.org/zap"
)

var (
	db     *mgo.Database
	session *mgo.Session
	logger *zap.SugaredLogger
)

func Init() {
	InitLogger()

	err := initDB()
	if err != nil {
		fmt.Println(err)
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

	session, err = mgo.Dial(dbConfig.ConnectionURL())
	fmt.Println(session)
	if err != nil {
		return
	}
	db = session.DB("function_junction")
	if err = session.Ping(); err != nil {
		return
	}

	// db.SetMaxIdleConns(dbConfig.MaxPoolSize())
	// db.SetMaxOpenConns(dbConfig.MaxOpenConns())
	// db.SetConnMaxLifetime(time.Duration(dbConfig.MaxLifeTimeMins()) * time.Minute)

	return
}

func GetDB() *mgo.Database {
	return db
}

func GetLogger() *zap.SugaredLogger {
	return logger
}

func Close() {
	logger.Sync()
	session.Close()
}
