package main

import (
	"pionira/cmd/api/handlers"
)

func (app *Application) routes(handler handlers.Handler) {

	apiGroup := app.server.Group("/api")
	publicAuthGroup := apiGroup.Group("/auth")
	{
		publicAuthGroup.POST("/register", handler.RegisterHandler)
		publicAuthGroup.POST("/login", handler.LoginHandler)
	}

	profileRoute := apiGroup.Group("/profile", app.AppMiddlewares.AuthenticationMiddleware)
	{
		profileRoute.GET("/authenticated/user", handler.GetAuthenticatedUser)
		profileRoute.PATCH("/change/password", handler.ChangeUserPassword)
	}
	app.server.GET("/", handler.HealthCheck)

}
