package main

import (
    "encoding/json"
    "net/http"
)

var testCampaignId1 = ""
var testCampaignId2 = ""

func testCampaignAdd1(client *http.Client) {
    data, code := transactClient(client, "POST", "http://localhost:8081/api/v1/campaign", "{\"name\":\"Test 1\",\"Description\":\"Madness\"}")
    test("testCampaignAdd1 code", code == 200)
    var result map[string]any
    json.Unmarshal(data, &result)

    test("campaign id", result["id"].(string) != "")
    test("campaign name", result["name"].(string) == "Test 1")
    test("campaign description", result["description"].(string) == "Madness")
    test("campaign creator", result["creator"].(string) == testCreatorId)
    test("campaign created", result["created"].(string) != "")
    test("campaign modified", result["modified"].(string) != "")
    test("campaign root", result["root"].(string) == "000000000000000000000000")
    test("campaign layout", result["layout"].(string) == "000000000000000000000000")

    testCampaignId1 = result["id"].(string)
}

func testCampaignAdd2(client *http.Client) {
    data, code := transactClient(client, "POST", "http://localhost:8081/api/v1/campaign", "{\"name\":\"Test 2\",\"Description\":\"WTF\"}")
    test("testCampaignAdd2 code", code == 200)
    var result map[string]any
    json.Unmarshal(data, &result)

    test("campaign id", result["id"].(string) != "")
    test("campaign name", result["name"].(string) == "Test 2")
    test("campaign description", result["description"].(string) == "WTF")
    test("campaign creator", result["creator"].(string) == testCreatorId)
    test("campaign created", result["created"].(string) != "")
    test("campaign modified", result["modified"].(string) != "")
    test("campaign root", result["root"].(string) == "000000000000000000000000")
    test("campaign layout", result["layout"].(string) == "000000000000000000000000")

    testCampaignId2 = result["id"].(string)
}

func testCampaignChange2(client *http.Client) {
    data, code := transactClient(client, "PUT", "http://localhost:8081/api/v1/campaign/" + testCampaignId2, "{\"name\":\"Test 2.1\",\"Description\":\"WTFXXX\"}")
    test("testCampaignChange2 code", code == 200)
    var result map[string]any
    json.Unmarshal(data, &result)

    test("campaign id", result["id"].(string) == testCampaignId2)
    test("campaign name", result["name"].(string) == "Test 2.1")
    test("campaign description", result["description"].(string) == "WTFXXX")
    test("campaign creator", result["creator"].(string) == testCreatorId)
    test("campaign created", result["created"].(string) != "")
    test("campaign modified", result["modified"].(string) != "")
    test("campaign root", result["root"].(string) == "000000000000000000000000")
    test("campaign layout", result["layout"].(string) == "000000000000000000000000")
}

func testCampaignGet2(client *http.Client) {
    data, code := transactClient(client, "GET", "http://localhost:8081/api/v1/campaign/" + testCampaignId2, "")
    test("testCampaignGet2 code", code == 200)
    var result map[string]any
    json.Unmarshal(data, &result)

    test("campaign id", result["id"].(string) == testCampaignId2)
    test("campaign name", result["name"].(string) == "Test 2.1")
    test("campaign description", result["description"].(string) == "WTFXXX")
    test("campaign creator", result["creator"].(string) == testCreatorId)
    test("campaign created", result["created"].(string) != "")
    test("campaign modified", result["modified"].(string) != "")
    test("campaign root", result["root"].(string) == "000000000000000000000000")
    test("campaign layout", result["layout"].(string) == "000000000000000000000000")
}

func testCampaignDelete2(client *http.Client) {
    _, code := transactClient(client, "DELETE", "http://localhost:8081/api/v1/campaign/" + testCampaignId2, "")
    test("testCampaignDelete2 code", code == 200)
}

func testCampaignList1(client *http.Client) {
    data, code := transactClient(client, "GET", "http://localhost:8081/api/v1/campaign", "")
    test("testCampaignList1 code", code == 200)
    var results []map[string]any
    json.Unmarshal(data, &results)
    test("campaign count", len(results) == 1)

    test("campaign id", results[0]["id"].(string) == testCampaignId1)
    test("campaign creator", results[0]["creator"].(string) == testCreatorId)
    test("campaign desc", results[0]["description"].(string) == "Madness")
    test("campaign name", results[0]["name"].(string) == "Test 1")
}

func testCampaignList2(client *http.Client) {
    data, code := transactClient(client, "GET", "http://localhost:8081/api/v1/campaign", "")
    test("testCampaignList2 code", code == 200)
    var results []map[string]any
    json.Unmarshal(data, &results)
    test("campaign count", len(results) == 2)

    u1Idx := 0
    u2Idx := 1
    if (results[0]["id"].(string) == testCampaignId2) {
        u1Idx = 1
        u2Idx = 0
    }

    test("campaign id", results[u1Idx]["id"].(string) == testCampaignId1)
    test("campaign creator", results[u1Idx]["creator"].(string) == testCreatorId)
    test("campaign desc", results[u1Idx]["description"].(string) == "Madness")
    test("campaign name", results[u1Idx]["name"].(string) == "Test 1")

    test("campaign id", results[u2Idx]["id"].(string) == testCampaignId2)
    test("campaign creator", results[u2Idx]["creator"].(string) == testCreatorId)
    test("campaign desc", results[u2Idx]["description"].(string) == "WTFXXX")
    test("campaign name", results[u2Idx]["name"].(string) == "Test 2.1")
}

func runCampaignTests(client *http.Client) {
    // Create a campaign
    testCampaignAdd1(client)

    // List campaigns
    testCampaignList1(client)

    // Create another campaign
    testCampaignAdd2(client)

    // Change a campaign
    testCampaignChange2(client)
    testCampaignGet2(client)

    // List campaigns
    testCampaignList2(client)

    // Delete a campaign
    testCampaignDelete2(client)

    // List campaigns
    testCampaignList1(client)

    // Recreate the deleted campaign for later
    testCampaignAdd2(client)
}
