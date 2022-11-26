package models

import (
    "bytes"
    "crypto/rand"
    "encoding/base64"
    "github.com/labstack/echo/v4"
    "github.com/zeebo/blake3"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "log"
    "net/http"
    "ostium/auth"
    "ostium/db"
    "time"
)

type UserPermission struct {
    Resource primitive.ObjectID `json:"resource" bson:"resource"`
    Realm int `json:"realm" bson:"realm"`
    Op int `json:"op" bson:"op"`
}

type UserResourceOpPair struct {
    Resource primitive.ObjectID
    Op int
}

type UserResourceOpPairs []UserResourceOpPair

type UserPermissionMap [auth.RealmMAX]UserResourceOpPairs

type User struct {
    Id *primitive.ObjectID `json:"id" bson:"_id,omitempty"`
    Email string `json:"email" bson:"email"`
    PassHash [32]byte `json:"-" bson:"passhash"`
    PassSalt [32]byte `json:"-" bson:"passsalt"`
    Password string `json:"password" bson:"-"`
    CookieSalt [32]byte `json:"-" bson:"cookiesalt"`
    Name string `json:"name" bson:"name"`
    Created time.Time `json:"created" bson:"created"`
    Modified time.Time `json:"modified" bson:"modified"`
    Permissions []UserPermission `json:"-" bson:"permissions"`
    perms UserPermissionMap `bson:"-"`
}

func checksum(in []byte) (ret [32]byte) {
    // Declare secret key
    // FIXME read the secret key from an ENV variable
    secret := bytes.Repeat([]byte("1"), 32)

    // Hash secret key and input data
    return blake3.Sum256(append(secret, in...))
}

func CreateUser(req *User) (user *User) {
    // Create object
    user = new(User)

    // Set details
    user.Email = req.Email
    user.Name = req.Name

    // Generate salt and set password hash
    rand.Read(user.PassSalt[:])
    rand.Read(user.CookieSalt[:])
    hashme := append(user.PassSalt[:], []byte(req.Password)...)
    user.PassHash = blake3.Sum256(hashme)

    // Set timestamps
    user.Created = time.Now()
    user.Modified = user.Created
    return user
}

func (this *User) Update(req *User) {
    // Update details
    if len(req.Email) != 0 {
        this.Email = req.Email
    }

    if len(req.Name) != 0 {
        this.Name = req.Name
    }

    // Generate salt and set password hash
    if len(req.Password) != 0 {
        rand.Read(this.PassSalt[:])
        rand.Read(this.CookieSalt[:])
        hashme := append(this.PassSalt[:], []byte(req.Password)...)
        this.PassHash = blake3.Sum256(hashme)
    }

    // Set timestamps
    this.Modified = time.Now()
}

func filterOp(src UserResourceOpPairs, id primitive.ObjectID, op int) bool {
    for _, el := range src {
        if el.Resource != id {
            continue
        }
        if el.Op == op || el.Op == auth.OpAdmin {
            return true
        }
    }
    return false
}

func (this *User) CheckOp(realm int, op int, id *primitive.ObjectID) (ret bool) {
    // Check wildcard
    var zeroId primitive.ObjectID

    // Try by realm first
    rmap := this.perms[realm]
    if filterOp(rmap, zeroId, op) {
        return true
    }
    if id != nil && filterOp(rmap, *id, op) {
        return true
    }

    // Try by wildcard realm
    rmap = this.perms[auth.RealmAll]
    if filterOp(rmap, zeroId, op) {
        return true
    }
    if id != nil && filterOp(rmap, *id, op) {
        return true
    }

    // No permission
    return false
}

func (this *User) AddPermission(c echo.Context, realm int, class int, id primitive.ObjectID) (ret bool) {
    this.Permissions = append(this.Permissions, UserPermission{
        id,
        realm,
        class})
    return true
}

func (this *User) Output(c echo.Context) {
    sumSlice := this.CookieSalt[:]
    sumSlice = append(sumSlice, this.Id[:]...)
    sum := checksum(sumSlice)
    raw := append(sum[:], sumSlice...)

    // Set the result cookie
    cookie := new(http.Cookie)
    cookie.Name = "auth"

    // Base 64 encode the cookie value
    cookie.Value = base64.StdEncoding.EncodeToString(raw)
    cookie.Expires = time.Now().Add(24 * time.Hour)
    c.SetCookie(cookie)
    return
}

func (this *User) CheckPassword(pass string) bool {
    hashme := append(this.PassSalt[:], []byte(pass)...)
    hash := blake3.Sum256(hashme)
    return this.PassHash == hash
}

func UserFromCookieWithoutSet(c echo.Context) *User {
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

    if len(hashedRaw) != (32 + 32 + 12) {
        log.Printf("Got wrong cookie")
        return nil
    }

    // Extract the expected hash and oid
    var hash [32]byte
    var salt [32]byte
    copy(hash[:], hashedRaw[:32])
    copy(salt[:], hashedRaw[32:64])

    oidSlice := hashedRaw[64:76]

    // Create the hasher to validate the content
    if checksum(hashedRaw[32:]) != hash {
        log.Printf("Got bad auth sum")
        return nil
    }

    // Load the user record
    var oid primitive.ObjectID
    copy(oid[:], oidSlice)

    var user User
    err = db.Get(&user, "user", oid)
    if err != nil {
        log.Printf("Failed to load user record")
        return nil
    }

    // Check to see if the cookie has been resalted
    if user.CookieSalt != salt {
        log.Printf("Cookie salt fails to match")
        return nil
    }

    // Distill perms
    var perms UserPermissionMap
    for _, p := range user.Permissions {
        perms[p.Realm] = append(perms[p.Realm], UserResourceOpPair{
            p.Resource, p.Op })
    }
    user.perms = perms

    return &user
}

func UserFromCookie(c echo.Context) *User {
    user := UserFromCookieWithoutSet(c)

    // Re-set cookie
    if user != nil {
        user.Output(c)
    }
    return user
}

func (this *User) ClearCookie() {
    // Reset the cookie salt
    rand.Read(this.CookieSalt[:])

    // Change it in the database
    err := db.Update("user", bson.D{{"$set", bson.D{{"cookiesalt", this.CookieSalt}}}}, *this.Id)
    if err != nil {
        panic(err)
    }
}
