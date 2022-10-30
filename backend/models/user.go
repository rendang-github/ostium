package models

import (
    "github.com/labstack/echo/v4"
    "github.com/zeebo/blake3"
    "log"
    "net/http"
    "time"
    "bytes"
    "encoding/base64"
)

/// auth.User
type User struct {
    CValue []byte
    Hash [32]byte
    Id int `json:"id"`
    Username string `json:"username"`
    Password string `json:"password"`
    PassHash [32]byte `json:"passhash"`
}

func checksum(in []byte) (ret [32]byte) {
    // Declare secret key
    // FIXME read the secret key from an ENV variable
    secret := bytes.Repeat([]byte("1"), 32)

    // Hash secret key and input data
    return blake3.Sum256(append(secret, in...))
}

func (this *User) CheckOp(c echo.Context, realm int, op int, id *string) (ret bool) {
    // FIXME implement
    return true
}

func (this *User) AddPermission(c echo.Context, realm int, class int, id *string) (ret bool) {
    // FIXME implement
    return true
}

func (this *User) Output(c echo.Context) {
    hashedRaw := append(this.Hash[:], this.CValue...)

    // Set the result cookie
    cookie := new(http.Cookie)
    cookie.Name = "auth"

    // Base 64 encode the cookie value
    cookie.Value = base64.StdEncoding.EncodeToString(hashedRaw)
    cookie.Expires = time.Now().Add(24 * time.Hour)
    c.SetCookie(cookie)
    return
}

func UserFromLocal(user string, pass string) (ret *User) {
    // Check credentials
    if (user != "peter" || pass != "foo") {
        return nil
    }

    // Create our state structure
    raw := []byte("bazuzu")

    // Create return structure
    ret = new(User)
    ret.CValue = raw
    ret.Hash = checksum(raw)
    return ret
}

func UserFromCookie(c echo.Context) (ret *User) {
    cookie, err := c.Cookie("auth")
    if err != nil {
        log.Printf("Error reading cookie")
        return nil
    }

    // Base64 decode the content
    hashedRaw, err := base64.StdEncoding.DecodeString(cookie.Value)
    if err != nil {
        log.Printf("Error base64 decoding cookie")
        return nil
    }

    if len(hashedRaw) < 32 {
        log.Printf("Got short cookie")
        return nil
    }

    // Extract the expected hash and raw content
    ret = new(User)
    ret.CValue = hashedRaw[32:]
    copy(ret.Hash[:], hashedRaw[:32])

    // Create the hasher to validate the content
    if checksum(ret.CValue) != ret.Hash {
        log.Printf("Got bad auth sum")
        return nil
    }

    return ret
}
