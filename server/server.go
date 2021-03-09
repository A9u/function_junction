package server

import (
	"fmt"
	"strconv"

	"github.com/A9u/function_junction/config"
	"github.com/rs/cors"
	"github.com/urfave/negroni"
)

func StartAPIServer() {
	port := config.AppPort()

	c := cors.New(cors.Options{Debug: true, AllowedHeaders: []string{"authorization", "content-type"}, AllowedMethods: []string{"GET", "POST", "DELETE"}})
	server := negroni.Classic()
	server.Use(c)
	dependencies, err := initDependencies()
	server.Use(negroni.HandlerFunc(AuthMiddleware))
	if err != nil {
		panic(err)
	}

	router := initRouter(dependencies)
	server.UseHandler(router)

	addr := fmt.Sprintf(":%s", strconv.Itoa(port))
	server.Run(addr)
}
