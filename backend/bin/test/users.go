package main

import (
    "encoding/json"
    "net/http"
)

var testOtherId = "";

func testUserList1(client *http.Client) {
    data := transactClient(client, "GET", "http://localhost:8081/api/v1/user", "")
    var results []map[string]any
    json.Unmarshal(data, &results)
    test("user count", len(results) == 1)
    test("user 1 id", results[0]["id"].(string) == testCreatorId)
    test("user 1 email", results[0]["email"].(string) == "testing@warlordsofbeer.com")
    test("user 1 name", results[0]["name"].(string) == "Test Admin")
    test("user 1 created", results[0]["created"].(string) != "")
    test("user 1 modified", results[0]["modified"].(string) != "")
}

func testUserAdd(client *http.Client) {
    data := transactClient(client, "POST", "http://localhost:8081/api/v1/user", "{\"email\":\"test2@warlordsofbeer.com\",\"password\":\"Test Password\",\"name\":\"Test User 2\"}")
    var result map[string]any
    json.Unmarshal(data, &result)

    test("new user id", result["id"].(string) != "")
    test("new user email", result["email"].(string) == "test2@warlordsofbeer.com")
    test("new user name", result["name"].(string) == "Test User 2")
    test("new user created", result["created"].(string) != "")
    test("new user modified", result["modified"].(string) != "")

    testOtherId = result["id"].(string)
}

func testUserChange(client *http.Client) {
    data := transactClient(client, "PUT", "http://localhost:8081/api/v1/user/" + testOtherId, "{\"name\":\"Doctor Rockso\"}")
    var result map[string]any
    json.Unmarshal(data, &result)

    test("new user id", result["id"].(string) == testOtherId)
    test("new user email", result["email"].(string) == "test2@warlordsofbeer.com")
    test("new user name", result["name"].(string) == "Doctor Rockso")
    test("new user created", result["created"].(string) != "")
    test("new user modified", result["modified"].(string) != "")

    testOtherId = result["id"].(string)
}

func testUserGet(client *http.Client) {
    data := transactClient(client, "GET", "http://localhost:8081/api/v1/user/" + testOtherId, "")
    var result map[string]any
    json.Unmarshal(data, &result)

    test("new user id", result["id"].(string) == testOtherId)
    test("new user email", result["email"].(string) == "test2@warlordsofbeer.com")
    test("new user name", result["name"].(string) == "Doctor Rockso")
    test("new user created", result["created"].(string) != "")
    test("new user modified", result["modified"].(string) != "")

    testOtherId = result["id"].(string)
}

func testUserDelete(client *http.Client) {
    transactClient(client, "DELETE", "http://localhost:8081/api/v1/user/" + testOtherId, "")
}

func testUserList2(client *http.Client) {
    data := transactClient(client, "GET", "http://localhost:8081/api/v1/user", "")
    var results []map[string]any
    json.Unmarshal(data, &results)
    test("user count", len(results) == 2)

    u1Idx := 0
    u2Idx := 1
    if (results[0]["id"].(string) != testCreatorId) {
        u1Idx = 1
        u2Idx = 0

    }
    test("orig user id", results[u1Idx]["id"].(string) == testCreatorId)
    test("orig user email", results[u1Idx]["email"].(string) == "testing@warlordsofbeer.com")
    test("orig user name", results[u1Idx]["name"].(string) == "Test Admin")

    test("new user id", results[u2Idx]["id"].(string) != testCreatorId)
    test("new user id", results[u2Idx]["id"].(string) != "")
    test("new user email", results[u2Idx]["email"].(string) == "test2@warlordsofbeer.com")
    test("new user name", results[u2Idx]["name"].(string) == "Doctor Rockso")
}

func runUserTests(client *http.Client) {
    // Get list of users
    testUserList1(client)

    // Create another user
    testUserAdd(client)

    // Change the user
    testUserChange(client)
    testUserGet(client)

    // Get list of users
    testUserList2(client)

    // Delete the new user and check that we get the original list
    testUserDelete(client)
    testUserList1(client)
}
