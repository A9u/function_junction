package server

import (
	"github.com/joshsoftware/golang-boilerplate/app"
	"github.com/joshsoftware/golang-boilerplate/category"
	"github.com/joshsoftware/golang-boilerplate/db"
)

type dependencies struct {
	CategoryService category.Service
}

func initDependencies() (dependencies, error) {
	appDB := app.GetDB()
	logger := app.GetLogger()
	dbStore := db.NewStorer(appDB)

	categoryService := category.NewService(dbStore, logger, app.GetCollection("catagories"))

	return dependencies{
		CategoryService: categoryService,
	}, nil
}
