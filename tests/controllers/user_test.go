package controllerTests

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/pick-up-api/router"
)

var (
	server   *httptest.Server
	reader   io.Reader //Ignore this for now
	usersUrl string
)

func init() {
	server = httptest.NewServer(router.GetRouter())
	usersUrl = fmt.Sprintf("%s/user", server.URL)
}

func TestUserProfile(t *testing.T) {
	validUserResponse, _ := http.Get(usersUrl + "/1")
	outOfBoundsResponse, _ := http.Get(usersUrl + "/11234123")
	defer validUserResponse.Body.Close()
	defer outOfBoundsResponse.Body.Close()

}

func TestCreateUser(t *testing.T) {
	emptyUser := url.Values{}

	userWithoutEmail := url.Values{}
	userWithoutEmail.Add("password", "secret")
	userWithoutEmail.Add("name", "Tester")

	userWithInvalidEmail := url.Values{}
	userWithInvalidEmail.Add("email", "test.test")
	userWithInvalidEmail.Add("password", "secret")
	userWithInvalidEmail.Add("name", "Tester")

	userWithoutPassword := url.Values{}
	userWithoutPassword.Add("email", "test.test")
	userWithoutPassword.Add("name", "Tester")

	user := url.Values{}
	user.Add("email", "test@yahoo.com")
	user.Add("password", "secret555")
	user.Add("name", "Tester")

	fails := []url.Values{emptyUser, userWithoutEmail, userWithInvalidEmail, userWithoutPassword}
	successes := []url.Values{user}

	// test invalid params
	for _, invalidUser := range fails {
		response, err := http.PostForm(usersUrl+"/create", invalidUser)
		defer response.Body.Close()

		errorIfExists(t, err)
		validateResponseCode(t, http.StatusBadRequest, response.StatusCode)
	}

	// test valid user json
	for _, validUser := range successes {
		response, err := http.PostForm(usersUrl+"/create", validUser)
		defer response.Body.Close()

		errorIfExists(t, err)
		validateResponseCode(t, http.StatusOK, response.StatusCode)

		responseBody, _ := ioutil.ReadAll(response.Body)
		var responseMap map[string]interface{}
		jsonErr := json.Unmarshal(responseBody, &responseMap)

		errorIfExists(t, jsonErr)
		code := responseMap["code"].(string)
		responseUserBody := responseMap["response"].(map[string]interface{})
		userId := responseUserBody["id"].(float64)
		userEmail := responseUserBody["email"].(string)
		userName := responseUserBody["name"].(string)

		validateCode(t, code)
		validateUserId(t, int64(userId))
		validateUserAttribute(t, userEmail, validUser.Get("email"), "email")
		validateUserPassword(t, responseUserBody["password"])
		validateUserAttribute(t, userName, validUser.Get("name"), "name")
	}
}

// func TestCreateUpdateAndDeleteUser(t *testing.T) {
//
// }

func errorIfExists(t *testing.T, err error) {
	if err != nil {
		t.Error(err)
	}
}

func validateResponseCode(t *testing.T, expectedStatusCode int, responseCode int) {
	if responseCode != expectedStatusCode {
		t.Errorf("Response status code expected %d. Received %d", expectedStatusCode, responseCode)
	}
}

func validateCode(t *testing.T, code string) {
	if code != "ok" {
		t.Errorf("Response code should be ok. Received %s")
	}
}

func validateUserId(t *testing.T, userId int64) {
	if int64(userId) <= 0 {
		t.Errorf("Response user id should be positive. Received %d", userId)
	}
}

func validateUserAttribute(t *testing.T, value string, expectedValue string, key string) {
	if value != expectedValue {
		t.Errorf("Response user %s changed from %s to %s", key, value, expectedValue)
	}
}

func validateUserPassword(t *testing.T, pw interface{}) {
	if pw != nil {
		t.Errorf("Response user password should not be returned. Received %s", pw)
	}
}
