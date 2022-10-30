package controllers
import (
    "github.com/labstack/echo/v4"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "net/http"
    "ostium/auth"
    "ostium/db"
    "ostium/models"
    "strings"
)

func APIUserPost(c echo.Context) (err error) {
    login := models.UserFromCookie(c)
    if login == nil || !login.CheckOp(c, auth.RealmUser, auth.OpCreate, nil) {
        return c.NoContent(http.StatusUnauthorized)
    }

    // Read the parameters from the POST
    reqUser := new(models.User)
    if err = c.Bind(reqUser); err != nil {
        return
    }

    // See if this user already exists
    var existUser models.User
    reqUser.Email = strings.ToLower(reqUser.Email)
    err = db.Query(&existUser, "user", bson.M{"email": reqUser.Email})
    if err == nil {
        // User already exists
        return c.NoContent(http.StatusConflict)
    }

    // Check we have a non-empty password
    if len(reqUser.Password) == 0 {
        return c.NoContent(http.StatusBadRequest)
    }

    // Create a new user
    user := models.CreateUser(reqUser)

    // Persist the record and get a new id
    oid := db.Insert("user", user)
    user.Id = &oid
    return c.JSON(http.StatusOK, user)
}

func APIUserPut(c echo.Context) (err error) {
    // Get the object id
    id := c.Param("id")
    oid, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return c.NoContent(http.StatusBadRequest)
    }

    login := models.UserFromCookie(c)
    oidHex := oid.Hex()
    if login == nil || !login.CheckOp(c, auth.RealmUser, auth.OpChange, &oidHex) {
        return c.NoContent(http.StatusUnauthorized)
    }

    // Read the parameters from the PUT
    reqUser := new(models.User)
    if err = c.Bind(reqUser); err != nil {
        return
    }

    // Read from the DB
    var existUser models.User
    err = db.Get(&existUser, "user", oid)
    if err != nil {
        return c.NoContent(http.StatusNotFound)
    }

    // See if this user already exists if we're changing the email
    reqUser.Email = strings.ToLower(reqUser.Email)
    if len(reqUser.Email) != 0 && existUser.Email != reqUser.Email {
        var existEmail models.User
        err = db.Query(&existEmail, "user", bson.M{"email": reqUser.Email})
        if err != nil {
            return c.NoContent(http.StatusConflict)
        }
    }

    // Update values
    existUser.Update(reqUser)

    // Persist the record
    err = db.Set("user", existUser, oid)
    if err != nil {
        return c.NoContent(http.StatusNotFound)
    }

    return c.JSON(http.StatusOK, existUser)
}

func APIUserGet(c echo.Context) (err error) {
    // Get the object id
    id := c.Param("id")
    oid, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return c.NoContent(http.StatusBadRequest)
    }

    // Check auth
    login := models.UserFromCookie(c)
    oidHex := oid.Hex()
    if login == nil || !login.CheckOp(c, auth.RealmUser, auth.OpRetrieve, &oidHex) {
        return c.NoContent(http.StatusUnauthorized)
    }

    // Read from the DB
    var user models.User
    err = db.Get(&user, "user", oid)
    if err != nil {
        return c.NoContent(http.StatusNotFound)
    }

    return c.JSON(http.StatusOK, user)
}

func APIUserAll(c echo.Context) (err error) {
    // Check auth
    login := models.UserFromCookie(c)
    if login == nil || !login.CheckOp(c, auth.RealmUser, auth.OpRetrieve, nil) {
        return c.NoContent(http.StatusUnauthorized)
    }

    // We want all records
    var users []models.User
    err = db.All(&users, "user")
    if err != nil {
        panic(err)
        return c.NoContent(http.StatusNotFound)
    }
    return c.JSON(http.StatusOK, users)
}

func APIUserDelete(c echo.Context) (err error) {
    // Get the object id
    id := c.Param("id")
    oid, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return c.NoContent(http.StatusBadRequest)
    }

    login := models.UserFromCookie(c)
    oidHex := oid.Hex()
    if login == nil || !login.CheckOp(c, auth.RealmUser, auth.OpErase, &oidHex) {
        return c.NoContent(http.StatusUnauthorized)
    }

    // Erase the record
    err = db.Delete("user", oid)
    if err != nil {
        return c.NoContent(http.StatusNotFound)
    }

    return c.NoContent(http.StatusOK)
}
