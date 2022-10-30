package main

import (
    "github.com/labstack/echo/v4"
    "github.com/labstack/echo/v4/middleware"
    "net/http"
    "ostium/controllers"
    "ostium/db"
)

func main() {
    // Set up database
    db.Initialize()

    // Set up web server
    server := echo.New()
    server.Use(middleware.Logger())
    server.Use(middleware.Recover())
    server.Use(middleware.CORS())
    server.Use(middleware.Secure())

    // Root route
    server.GET("/", func(c echo.Context) error {
        return c.String(http.StatusOK, "Hello, World!")
    })

    // Login routes
    server.POST("/api/v1/login", controllers.APILoginPost)
    server.GET("/api/v1/login", controllers.APILoginGet)

    // User routes
    server.POST("/api/v1/user", controllers.APIUserPost)
    server.PUT("/api/v1/user/:id", controllers.APIUserPut)
    server.GET("/api/v1/user", controllers.APIUserAll)
    server.GET("/api/v1/user/:id", controllers.APIUserGet)
    server.DELETE("/api/v1/user/:id", controllers.APIUserDelete)

    // Campaign routes
    server.POST("/api/v1/campaign", controllers.APICampaignPost)
    server.PUT("/api/v1/campaign/:id", controllers.APICampaignPut)
    server.GET("/api/v1/campaign", controllers.APICampaignAll)
    server.GET("/api/v1/campaign/:id", controllers.APICampaignGet)
    server.DELETE("/api/v1/campaign/:id", controllers.APICampaignDelete)

    server.File("/", "public/index.html")
    server.File("/index.html", "public/index.html")
    server.File("/favicon.png", "public/favicon.png")

    // Launch web server
    server.Logger.Fatal(server.Start(":8081"))
}
