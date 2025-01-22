package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"os"
	"pionira/cmd/api/handlers"
	"pionira/cmd/api/middlewares"
	"pionira/common"
	"pionira/internal/mailer"
)

type Application struct {
	logger         echo.Logger
	server         *echo.Echo
	handler        handlers.Handler
	AppMiddlewares middlewares.AppMiddleware
}

func main() {
	e := echo.New()
	e.Use(middlewares.CustomMiddleware)
	e.Use(middlewares.AnotherMiddleware)
	e.Use(middleware.Logger())

	err := godotenv.Load()
	if err != nil {
		e.Logger.Fatal(err.Error())
	}

	db, err := common.NewMysql()
	if err != nil {
		e.Logger.Fatal("Error loading .env file")
	}

	appMailer := mailer.NewMailer(e.Logger)

	h := handlers.Handler{
		DB:     db,
		Logger: e.Logger,
		Mailer: appMailer,
	}

	appMiddleware := middlewares.AppMiddleware{
		Logger: e.Logger,
		DB:     db,
	}

	app := Application{
		logger:         e.Logger,
		server:         e,
		handler:        h,
		AppMiddlewares: appMiddleware,
	}

	//Route
	app.routes(h)
	fmt.Println(app)

	port := os.Getenv("PORT")
	addr := fmt.Sprintf("localhost:%s", port)
	e.Logger.Fatal(e.Start(addr))

}
