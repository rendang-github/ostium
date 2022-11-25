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

    // Test the user controller
    runUserTests(client)

    // Test the user controller
    runCampaignTests(client)

    log.Printf("%d/%d tests succeeded", testNumber - testFails, testNumber)
    //time.Sleep(60 * 60 * 24 * time.Second)
}
