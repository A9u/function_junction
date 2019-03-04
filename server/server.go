package server

import (
	"fmt"
	"strconv"

	"github.com/A9u/function_junction/config"
	"github.com/urfave/negroni"
)

func StartAPIServer() {
	port := config.AppPort()

	server := negroni.Classic()

	server.Use(negroni.HandlerFunc(AuthMiddleware))

	dependencies, err := initDependencies()
	if err != nil {
		panic(err)
	}

	router := initRouter(dependencies)
	server.UseHandler(router)

	addr := fmt.Sprintf(":%s", strconv.Itoa(port))
	server.Run(addr)
}
