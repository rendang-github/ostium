package auth

// Permissions realm
const (
    RealmAll = 0
    RealmUser = 1
    RealmCampaign = 2
    RealmLayout = 3
    RealmSnippet = 4
    RealmResource = 5
    RealmMAX = 6
)

// Permissions operation
const (
    OpCreate = 1
    OpRetrieve = 2
    OpChange = 3
    OpErase = 4
    OpAdmin = 5
)

