package main

import (
    "encoding/json"
    "net/http"
    "strconv"
)

func testPermissionCampaignGet(client *http.Client, id string, want int) bool {
    data, code := transactClient(client, "GET", "http://localhost:8081/api/v1/campaign/" + id, "")
    test("testPermissionCampaignGet code", code == want)
    var result map[string]any
    json.Unmarshal(data, &result)
    return len(result) != 0 && want == code
}

func testPermissionCampaignSet(client *http.Client, id string, want int) bool {
    data, code := transactClient(client, "PUT", "http://localhost:8081/api/v1/campaign/" + id, "{\"name\":\"Moo\",\"Description\":\"Cow\"}")
    test("testPermissionCampaignSet code", code == want)
    var result map[string]any
    json.Unmarshal(data, &result)
    if len(result) == 0 {
        return false
    }

    test("perm campaign id", result["id"].(string) == id)
    test("perm campaign name", result["name"].(string) == "Moo")
    test("perm campaign description", result["description"].(string) == "Cow")
    return true
}

func testPermissionCampaignDelete(client *http.Client, id string, want int) bool {
    _, code := transactClient(client, "DELETE", "http://localhost:8081/api/v1/campaign/" + id, "")
    test("testPermissionCampaignDelete code", code == want)
    return code == 200
}

func testAddPermission(client *http.Client, userId string, objectId string, realm int, op int) {
    data, code := transactClient(client, "POST", "http://localhost:8081/api/v1/permissions/" + userId, "{\"add\":[{\"resource\":\"" + objectId + "\",\"realm\":" + strconv.Itoa(realm) + ",\"op\":" + strconv.Itoa(op) + "}]}")
    test("testAddPermission code", code == 200)
    var result map[string]any
    json.Unmarshal(data, &result)

    test("perm add user id", result["id"].(string) == userId)
}

func testRemovePermission(client *http.Client, userId string, objectId string, realm int, op int) {
    data, code := transactClient(client, "POST", "http://localhost:8081/api/v1/permissions/" + userId, "{\"remove\":[{\"resource\":\"" + objectId + "\",\"realm\":" + strconv.Itoa(realm) + ",\"op\":" + strconv.Itoa(op) + "}]}")
    test("testAddPermission code", code == 200)
    var result map[string]any
    json.Unmarshal(data, &result)

    test("perm remove user id", result["id"].(string) == userId)
}

func runPermissionTests(mainClient *http.Client, altClient *http.Client, altId string) {
    // Test that the alt-client can't access any of the campaigns we created
    test("initial campaign 1", !testPermissionCampaignGet(altClient, testCampaignId1, 401))
    test("initial campaign 2", !testPermissionCampaignGet(altClient, testCampaignId2, 401))

    // Add permission to read the first campaign
    testAddPermission(mainClient, altId, testCampaignId1, 2, 2)

    // Test that the alt-client can access only the first campaign
    test("perm 1 campaign 1", testPermissionCampaignGet(altClient, testCampaignId1, 200))
    test("perm 1 campaign 2", !testPermissionCampaignGet(altClient, testCampaignId2, 401))

    // Add permission to read the second campaign
    testAddPermission(mainClient, altId, testCampaignId2, 2, 2)
    test("perm 2 campaign 1", testPermissionCampaignGet(altClient, testCampaignId1, 200))
    test("perm 2 campaign 2", testPermissionCampaignGet(altClient, testCampaignId2, 200))

    // Revoke permission to read the first campaign
    testRemovePermission(mainClient, altId, testCampaignId1, 2, 2)
    test("perm 3 campaign 1", !testPermissionCampaignGet(altClient, testCampaignId1, 401))
    test("perm 3 campaign 2", testPermissionCampaignGet(altClient, testCampaignId2, 200))

    // Test write permission
    test("perm 4 campaign 1", !testPermissionCampaignSet(altClient, testCampaignId1, 401))
    test("perm 4 campaign 2", !testPermissionCampaignSet(altClient, testCampaignId2, 401))
    testAddPermission(mainClient, altId, testCampaignId2, 2, 3)
    test("perm 5 campaign 1", !testPermissionCampaignSet(altClient, testCampaignId1, 401))
    test("perm 5 campaign 2", testPermissionCampaignSet(altClient, testCampaignId2, 200))

    // Test delete permission
    test("perm 6 campaign 1", !testPermissionCampaignDelete(altClient, testCampaignId1, 401))
    test("perm 6 campaign 2", !testPermissionCampaignDelete(altClient, testCampaignId2, 401))
    testAddPermission(mainClient, altId, testCampaignId2, 2, 4)
    test("perm 7 campaign 1", !testPermissionCampaignDelete(altClient, testCampaignId1, 401))
    test("perm 7 campaign 2", testPermissionCampaignDelete(altClient, testCampaignId2, 200))
    test("perm 8 campaign 2", !testPermissionCampaignGet(altClient, testCampaignId2, 404))

}
