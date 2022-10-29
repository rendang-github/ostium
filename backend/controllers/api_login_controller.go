package controllers

import (
    "ostium/models"
    "github.com/labstack/echo/v4"
    "net/http"
)

/// Login request
type RequestLogin struct {
    // Requested username
    Username string `json:"username"`

    // Requested password
    Password string `json:"password"`
}

/// Login response
type ResponseLogin struct {
    // Success flag
    Success bool `json:"success"`
}

func APILoginPost(c echo.Context) (err error) {
    req := new(RequestLogin)
    if err = c.Bind(req); err != nil {
        return
    }

    // Check the password
    //id := models.SetExit(e)

    // Check the values
    user := models.UserFromLocal(req.Username, req.Password)
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
