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
		publicAuthGroup.POST("/forgot/password", handler.ForgotPasswordHandler)
		publicAuthGroup.POST("/reset/password", handler.ResetPasswordHandler)
	}

	profileRoute := apiGroup.Group("/profile", app.AppMiddlewares.AuthenticationMiddleware)
	{
		profileRoute.GET("/authenticated/user", handler.GetAuthenticatedUser)
		profileRoute.PATCH("/change/password", handler.ChangeUserPassword)
	}

	categoryRoute := apiGroup.Group("/category", app.AppMiddlewares.AuthenticationMiddleware)
	{
		categoryRoute.GET("/all", handler.ListCategories)
	}
	app.server.GET("/", handler.HealthCheck)

}
