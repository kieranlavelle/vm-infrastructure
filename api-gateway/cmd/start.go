package main

import (
	"github.com/kieranlavelle/api_gateway/pkg/api"
	"github.com/kieranlavelle/api_gateway/pkg/tokenserver"
)

func main() {

	// create our routers for the apis
	apiRouter := api.CreateRoutes()
	tokenServerRouter := tokenserver.CreateTokenServer()

	go apiRouter.Run("0.0.0.0:8002")
	tokenServerRouter.Run("0.0.0.0:8006")
}
