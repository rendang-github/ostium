package main

import (
    "crypto/tls"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "net/http/cookiejar"
    "os/exec"
    "ostium/config"
    "ostium/db"
    "ostium/models"
    "strings"
    "time"
)

var seedUserName = "testing@warlordsofbeer.com"
var seedUserPass = "My Milkshake"

func createClient() (*http.Client) {
    cookieJar, _ := cookiejar.New(nil)
    tr := &http.Transport{
        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
    }

    client := &http.Client{Transport: tr, Jar: cookieJar, CheckRedirect: func(req *http.Request, via []*http.Request) error {
        return http.ErrUseLastResponse
    }}

    return client
}


func transactClient(client *http.Client, method string, url string, requestBody string) []byte {
    req, err := http.NewRequest(method, url, strings.NewReader(requestBody))
    req.Header.Set("Content-Type", "application/json")
    if err != nil {
        fmt.Printf("client: could not create request: %s\n", err)
        var ret []byte
        return ret;
    }
    res, err := client.Do(req)
    if err != nil {
        fmt.Printf("%s", err)
    }
    defer res.Body.Close()

    fmt.Println(res.Header)
    responseBody, err := ioutil.ReadAll(res.Body)
    if err != nil {
        fmt.Printf("%s", err)
    }
    fmt.Println(string(responseBody))
    return responseBody
}

var testNumber = 0;
var testFails = 0;
var testCreatorId = "";
func test(label string, success bool) {
    testNumber++
    if (!success) {
        testFails++
        log.Printf("Test %d \"%s\" failed", testNumber, label)
    }
}

func testLogin(client *http.Client) {
    data := transactClient(client, "POST", "http://localhost:8081/api/v1/login", "{\"username\":\"testing@warlordsofbeer.com\",\"password\":\"My Milkshake\"}")
    var result map[string]any
    json.Unmarshal(data, &result)

    test("login id", result["id"].(string) != "")
    test("login email", result["email"].(string) == "testing@warlordsofbeer.com")
    test("login name", result["name"].(string) == "Test Admin")
    test("login created", result["created"].(string) != "")
    test("login modified", result["modified"].(string) != "")

    testCreatorId = result["id"].(string)
}

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
    test("new user name", results[u2Idx]["name"].(string) == "Test User 2")
}

func testCampaignAdd(client *http.Client) {
    data := transactClient(client, "POST", "http://localhost:8081/api/v1/campaign", "{\"name\":\"Test 1\",\"Description\":\"Madness\"}")
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
}

func testCampaignList(client *http.Client) {
    data := transactClient(client, "GET", "http://localhost:8081/api/v1/campaign", "")
    var results []map[string]any
    json.Unmarshal(data, &results)
    test("campaign count", len(results) == 1)

    test("campaign id", results[0]["id"].(string) != "")
    test("campaign creator", results[0]["creator"].(string) == testCreatorId)
    test("campaign desc", results[0]["description"].(string) == "Madness")
    test("campaign name", results[0]["name"].(string) == "Test 1")
}

func main() {
    config.DatabaseURI = "mongodb://127.0.0.1:27017"
    config.DatabaseName = "ostium_test"

    // Clear out the test database
    db.Clear()

    // Start up the web server
    cmd := exec.Command("./ostium", "-dbname", "ostium_test")
    err := cmd.Start()
    defer cmd.Process.Kill()
    if err != nil {
        log.Fatal(err)
    }

    // Create a seed user
    seedUser := new(models.User)
    seedUser.Email = seedUserName
    seedUser.Name = "Test Admin"
    seedUser.Password = seedUserPass
    saveUser := models.CreateUser(seedUser);
    db.Insert("user", saveUser)

    // Give the web server a second to start up
    time.Sleep(1 * time.Second)

    client := createClient()

    // Attempt to log in
    testLogin(client)

    // Get list of users
    testUserList1(client)

    // Create another user
    testUserAdd(client)

    // Get list of users
    testUserList2(client)

    // Create a campaign
    testCampaignAdd(client)

    // List campaigns
    testCampaignList(client)

    log.Printf("%d/%d tests succeeded", testNumber - testFails, testNumber)
    //time.Sleep(60 * 60 * 24 * time.Second)
}
