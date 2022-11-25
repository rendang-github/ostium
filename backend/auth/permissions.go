package auth

// Permissions realm
const (
    RealmUser = 1
    RealmCampaign = 2
    RealmLayout = 3
    RealmSnippet = 4
    RealmResource = 5
)

// Permissions operation
const (
    OpCreate = 1
    OpRetrieve = 2
    OpChange = 3
    OpErase = 4
)

// Ownership class
const (
    ClassOwner = 1
    ClassContributor = 2
    ClassViewer = 3
)
