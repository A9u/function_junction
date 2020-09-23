package server

import (
	"fmt"
	"strconv"

	"github.com/joshsoftware/function_junction/config"
	"github.com/urfave/negroni"
)

func StartAPIServer() {
	port := config.AppPort()

	server := negroni.Classic()
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
