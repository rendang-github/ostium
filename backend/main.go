package main

import (
    "net/http"
    "github.com/labstack/echo/v4"
    "github.com/labstack/echo/v4/middleware"
    "ostium/controllers"
)

func main() {
    server := echo.New()
    server.Use(middleware.Logger())
    server.Use(middleware.Recover())
    server.Use(middleware.CORS())
    server.Use(middleware.Secure())

    server.GET("/", func(c echo.Context) error {
        return c.String(http.StatusOK, "Hello, World!")
    })
    server.POST("/api/v1/login", controllers.APILoginPost)
    server.GET("/api/v1/test", controllers.APITestGet)

    server.File("/", "public/index.html")
    server.File("/index.html", "public/index.html")
    server.File("/favicon.png", "public/favicon.png")

    server.Logger.Fatal(server.Start(":8081"))
}
