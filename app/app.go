package app

import (
	"time"

	"github.com/go-mgo/mgo"
	"github.com/joshsoftware/golang-boilerplate/config"
	"go.uber.org/zap"
)

var (
	db     *mgo.Session
	logger *zap.SugaredLogger
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
	db, err = mgo.Dial(dbConfig.ConnectionURL())
	if err != nil {
	  return
	}

	if err = db.Ping(); err != nil {
	  return
	}

	// db.SetMaxIdleConns(dbConfig.MaxPoolSize())
	// db.SetMaxOpenConns(dbConfig.MaxOpenConns())
	// db.SetConnMaxLifetime(time.Duration(dbConfig.MaxLifeTimeMins()) * time.Minute)

	return
}

func GetDB() *mgo.Session {
 	return db
}

func GetLogger() *zap.SugaredLogger {
	return logger
}

func Close() {
	logger.Sync()
		db.Close()
}

func GetCollection(dbname string, col string) *mgo.Collection {
  return db.DB(dbname).C(col)
}

func Copy() *mgo.Session {
  return &db.Copy()
}

// func(s *Session) Close() {
//   if(s.session != nil) {
//     s.session.Close()
//   }
// }
