package controllers
import (
    "github.com/labstack/echo/v4"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "net/http"
    "ostium/models"
)

func Check(c echo.Context, realm int, op int) (oid primitive.ObjectID, login *models.User, success bool) {
    // Read the user
    login = models.UserFromCookie(c)
    if login == nil {
        c.NoContent(http.StatusUnauthorized)
        return oid, login, false
    }

    // Check to see if there's meant to be an ID param
    names := c.ParamNames()
    found := false
    for _, label := range names {
        if label == "id" {
            found = true
            break
        }
    }

    // No object ID, just do a unscoped permissions check
    if found == false {
        if !login.CheckOp(realm, op, nil) {
            c.NoContent(http.StatusUnauthorized)
            return oid, login, false
        }
        return oid, login, true
    }

    // Get the object id
    id := c.Param("id")

    // Read oid
    oid, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        c.NoContent(http.StatusBadRequest)
        return oid, login, false
    }

    // Check auth
    if login == nil || !login.CheckOp(realm, op, &oid) {
        c.NoContent(http.StatusUnauthorized)
        return oid, login, false
    }

    return oid, login, true
}

