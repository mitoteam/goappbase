package goappbase

import (
	"log"

	"github.com/gin-gonic/gin"
)

func (app *AppBase) buildWebRouter() {
	//no debug logging
	gin.SetMode(gin.ReleaseMode)

	// Prepare router
	app.webRouter = gin.New()

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	app.webRouter.Use(gin.Recovery())

	if app.WebRouterLogQueries {
		app.webRouter.Use(gin.Logger())
		log.Println("Extended queries logging enabled.")
	}

	if app.BuildWebRouterF != nil {
		app.BuildWebRouterF(app.webRouter)
	}
}
