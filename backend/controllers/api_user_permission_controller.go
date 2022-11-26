package controllers
import (
    "github.com/labstack/echo/v4"
    "net/http"
    "ostium/auth"
    "ostium/db"
    "ostium/models"
    "sort"
)

/// Permissions change request
type PermissionsChangeRequest struct {
    // Permission objects to remove
    Remove []models.UserPermission `json:"remove"`

    // Permission objects to add
    Add []models.UserPermission `json:"add"`
}

func APIUserPermissionsPost(c echo.Context) (err error) {
    oid, _, success := Check(c, auth.RealmUser, auth.OpAdmin)
    if !success {
        return
    }

    // Read the parameters from the POST
    reqChange := new(PermissionsChangeRequest)
    if err = c.Bind(reqChange); err != nil {
        return
    }

    // Read from the DB
    var user models.User
    err = db.Get(&user, "user", oid)
    if err != nil {
        return c.NoContent(http.StatusNotFound)
    }

    // Remove requested permission objects
    for _, remove := range reqChange.Remove {
        for idx, p := range user.Permissions {
            if p == remove {
                // Remove efficiently by swapping remove element with last
                // element, then truncating the slice
                lastIdx := len(user.Permissions) - 1
                user.Permissions[idx] = user.Permissions[lastIdx]
                user.Permissions = user.Permissions[:lastIdx]
            }
        }
    }

    // Add new permission objects
    for _, add := range reqChange.Add {
        user.Permissions = append(user.Permissions, add)
    }

    // Deduplicate permissions, start with sorting the elements
    sort.Slice(user.Permissions, func(i, j int) bool {
        if (user.Permissions[i].Realm == user.Permissions[j].Realm) {
            if (user.Permissions[i].Resource == user.Permissions[j].Resource) {
                return user.Permissions[i].Op < user.Permissions[j].Op
            }
            return user.Permissions[i].Resource.Hex() < user.Permissions[j].Resource.Hex()
        }
        return user.Permissions[i].Realm < user.Permissions[j].Realm
    })
    endIdx := len(user.Permissions)
    writeIdx := 1
    for readIdx := 1; readIdx < endIdx; readIdx++ {
        if user.Permissions[readIdx] == user.Permissions[readIdx - 1] {
            // Same as last one, skip the write
            continue
        }
        if readIdx != writeIdx {
            // Shift the element down
            user.Permissions[writeIdx] = user.Permissions[readIdx]
        }
        writeIdx++
    }

    // Truncate leftovers
    user.Permissions = user.Permissions[:writeIdx]

    // Persist the record
    err = db.Set("user", user, oid)
    if err != nil {
        return c.NoContent(http.StatusNotFound)
    }

    return c.JSON(http.StatusOK, user)
}

func APIUserPermissionsGet(c echo.Context) (err error) {
    oid, _, success := Check(c, auth.RealmUser, auth.OpAdmin)
    if !success {
        return
    }

    // Read from the DB
    var user models.User
    err = db.Get(&user, "user", oid)
    if err != nil {
        return c.NoContent(http.StatusNotFound)
    }

    return c.JSON(http.StatusOK, user)
}
