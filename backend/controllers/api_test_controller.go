package controllers

import (
    "ostium/models"
    "github.com/labstack/echo/v4"
    "net/http"
)

func APITestGet(c echo.Context) (err error) {
    // Check the values
    user := models.UserFromCookie(c)
    if user == nil {
        return c.JSON(http.StatusOK, ResponseLogin {
            Success: false,
        })
    }

    // Set the result cookie
    user.Output(c)

    // Return a response
    return c.JSON(http.StatusOK, ResponseLogin {
        Success: true,
    })
}
