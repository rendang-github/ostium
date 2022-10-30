package controllers
import (
    "github.com/labstack/echo/v4"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "net/http"
    "ostium/auth"
    "ostium/db"
    "ostium/models"
    "time"
)

func APICampaignPost(c echo.Context) (err error) {
    login := models.UserFromCookie(c)
    if login == nil || !login.CheckOp(c, auth.RealmCampaign, auth.OpCreate, nil) {
        return c.NoContent(http.StatusUnauthorized)
    }

    // Read the parameters from the POST
    campaign := new(models.Campaign)
    if err = c.Bind(campaign); err != nil {
        return
    }

    // Set starting timestamps
    campaign.Created = time.Now()
    campaign.Modified = campaign.Created

    // Persist the record and get a new id
    oid := db.Insert("campaign", campaign)
    campaign.Id = &oid

    // Allow the current user to own the new campaign id
    hexId := oid.Hex();
    if ! login.AddPermission(c, auth.RealmCampaign, auth.ClassOwner, &hexId) {
        return c.NoContent(http.StatusInternalServerError)
    }

    return c.JSON(http.StatusOK, campaign)
}

func APICampaignPut(c echo.Context) (err error) {
    // Get the object id
    id := c.Param("id")
    oid, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return c.NoContent(http.StatusBadRequest)
    }

    login := models.UserFromCookie(c)
    oidHex := oid.Hex()
    if login == nil || !login.CheckOp(c, auth.RealmCampaign, auth.OpChange, &oidHex) {
        return c.NoContent(http.StatusUnauthorized)
    }

    // Read the parameters from the PUT
    newCampaign := new(models.Campaign)
    if err = c.Bind(newCampaign); err != nil {
        return
    }

    // Read from the DB
    var existCampaign models.Campaign
    err = db.Get(&existCampaign, "campaign", oid)
    if err != nil {
        return c.NoContent(http.StatusNotFound)
    }

    // Update values
    existCampaign.Name = newCampaign.Name
    existCampaign.Description = newCampaign.Description
    existCampaign.Modified = time.Now()

    // Persist the record
    err = db.Set("campaign", existCampaign, oid)
    if err != nil {
        return c.NoContent(http.StatusNotFound)
    }

    return c.JSON(http.StatusOK, existCampaign)
}

func APICampaignGet(c echo.Context) (err error) {
    // Get the object id
    id := c.Param("id")
    oid, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return c.NoContent(http.StatusBadRequest)
    }

    // Check auth
    login := models.UserFromCookie(c)
    oidHex := oid.Hex()
    if login == nil || !login.CheckOp(c, auth.RealmCampaign, auth.OpRetrieve, &oidHex) {
        return c.NoContent(http.StatusUnauthorized)
    }

    // Read from the DB
    var campaign models.Campaign
    err = db.Get(&campaign, "campaign", oid)
    if err != nil {
        return c.NoContent(http.StatusNotFound)
    }

    return c.JSON(http.StatusOK, campaign)
}

func APICampaignAll(c echo.Context) (err error) {
    // Check auth
    login := models.UserFromCookie(c)
    if login == nil || !login.CheckOp(c, auth.RealmCampaign, auth.OpRetrieve, nil) {
        return c.NoContent(http.StatusUnauthorized)
    }

    // We want all records
    var campaigns []models.Campaign
    err = db.All(&campaigns, "campaign")
    if err != nil {
        panic(err)
        return c.NoContent(http.StatusNotFound)
    }
    return c.JSON(http.StatusOK, campaigns)
}

func APICampaignDelete(c echo.Context) (err error) {
    // Get the object id
    id := c.Param("id")
    oid, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return c.NoContent(http.StatusBadRequest)
    }

    login := models.UserFromCookie(c)
    oidHex := oid.Hex()
    if login == nil || !login.CheckOp(c, auth.RealmCampaign, auth.OpErase, &oidHex) {
        return c.NoContent(http.StatusUnauthorized)
    }

    // Erase the record
    err = db.Delete("campaign", oid)
    if err != nil {
        return c.NoContent(http.StatusNotFound)
    }

    return c.NoContent(http.StatusOK)
}
