package server

import (
	"github.com/A9u/function_junction/app"
	"github.com/A9u/function_junction/category"
	"github.com/A9u/function_junction/db"
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
