package main

import (
    "crypto/tls"
    "encoding/json"
    "fmt"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "io"
    "io/ioutil"
    "log"
    "net/http"
    "net/http/cookiejar"
    "os"
    "os/exec"
    "ostium/auth"
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


func transactClient(client *http.Client, method string, url string, requestBody string) ([]byte, int) {
    req, err := http.NewRequest(method, url, strings.NewReader(requestBody))
    req.Header.Set("Content-Type", "application/json")
    if err != nil {
        fmt.Printf("client: could not create request: %s\n", err)
        var ret []byte
        return ret, 0;
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
    return responseBody, res.StatusCode
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

func logInClient(client *http.Client, username string, password string, idTarget *string) {
    data, code := transactClient(client, "POST", "http://localhost:8081/api/v1/login", "{\"username\":\"" + username + "\",\"password\":\"" + password + "\"}")
    test("logInClient code", code == 200)
    var result map[string]any
    json.Unmarshal(data, &result)

    if (idTarget != nil) {
        test("login id", result["id"].(string) != "")
        test("login email", result["email"].(string) == username)
        test("login created", result["created"].(string) != "")
        test("login modified", result["modified"].(string) != "")

        *idTarget = result["id"].(string)
    } else {
        // We weren't supposed to get a result
        test("No result", len(result) == 0)
    }
}

func copyToStdout(stream *io.ReadCloser) {
    io.Copy(os.Stdout, *stream)
}

func main() {
    config.DatabaseURI = "mongodb://127.0.0.1:27017"
    config.DatabaseName = "ostium_test"

    // Clear out the test database
    db.Clear()

    // Start up the web server
    cmd := exec.Command("./ostium", "-dbname", "ostium_test")
    stdout, err := cmd.StdoutPipe()
    if err != nil {
        log.Fatal(err)
    }
    stderr, err := cmd.StderrPipe()
    if err != nil {
        log.Fatal(err)
    }
    err = cmd.Start()
    defer cmd.Process.Kill()
    if err != nil {
        log.Fatal(err)
    }

    // Echo stdout, stderr in another goroutine
    go copyToStdout(&stdout)
    go copyToStdout(&stderr)

    // Create a seed user
    seedUser := new(models.User)
    seedUser.Email = seedUserName
    seedUser.Name = "Test Admin"
    seedUser.Password = seedUserPass
    saveUser := models.CreateUser(seedUser);
    saveUser.Permissions = append(saveUser.Permissions, models.UserPermission{
        primitive.ObjectID{},
        auth.RealmAll,
        auth.OpAdmin })
    db.Insert("user", saveUser)

    // Give the web server a second to start up
    time.Sleep(1 * time.Second)

    client := createClient()

    // Attempt to log in the main user
    logInClient(client, "testing@warlordsofbeer.com", "My Milkshake", &testCreatorId)

    // Test the user controller
    runUserTests(client)

    // Test the user controller
    runCampaignTests(client)

    // Log in the client we created in the user tests
    altclient := createClient()
    var checkId string
    logInClient(altclient, "test2@warlordsofbeer.com", "Test Password", &checkId)

    // Run permissions tests
    runPermissionTests(client, altclient, checkId)

    log.Printf("%d/%d tests succeeded", testNumber - testFails, testNumber)
    //time.Sleep(60 * 60 * 24 * time.Second)
}
