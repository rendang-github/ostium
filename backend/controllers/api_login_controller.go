package controllers

import (
    "github.com/labstack/echo/v4"
    "go.mongodb.org/mongo-driver/bson"
    "net/http"
    "ostium/db"
    "ostium/models"
    "strings"
)

/// Login request
type RequestLogin struct {
    // Requested username
    Username string `json:"username"`

    // Requested password
    Password string `json:"password"`
}

func APILoginGet(c echo.Context) (err error) {
    // Check for login
    user := models.UserFromCookie(c)
    if user == nil {
        return c.JSON(http.StatusUnauthorized, nil)
    }

    // Return a response
    return c.JSON(http.StatusOK, user)
}

func APILoginDelete(c echo.Context) (err error) {
    // Check for login
    user := models.UserFromCookieWithoutSet(c)
    if user == nil {
        return c.JSON(http.StatusOK, true)
    }

    // Clear the cookie
    user.ClearCookie()

    // Return a response
    return c.JSON(http.StatusOK, true)
}

func APILoginPost(c echo.Context) (err error) {
    req := new(RequestLogin)
    if err = c.Bind(req); err != nil {
        return
    }

    // Attempt to load the appropriate user record
    req.Username = strings.ToLower(req.Username)
    var user models.User
    err = db.Query(&user, "user", bson.M{"email": req.Username})
    if err != nil {
        return c.JSON(http.StatusOK, false)
    }

    // Check the password
    if ! user.CheckPassword(req.Password) {
        return c.JSON(http.StatusOK, false)
    }

    // Set the result cookie
    user.Output(c)

    // Return a response
    return c.JSON(http.StatusOK, user)
}
