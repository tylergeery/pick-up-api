package controllerTests

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/pick-up-api/router"
)

type user struct {
	id         int64
	email      string
	name       string
	token      string
	is_active  int
	created_at string
	updated_at string
}

type response struct {
	code string
	user user
}

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

	validateResponseCode(t, http.StatusOK, validUserResponse.StatusCode)
	validateResponseCode(t, http.StatusNotFound, outOfBoundsResponse.StatusCode)

	_ = validateGetUser(t, validUserResponse.Body, url.Values{})
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

		_ = validateGetUser(t, response.Body, user)
	}
}

func TestCreateUpdateAndDeleteUser(t *testing.T) {
	user := url.Values{}
	user.Add("email", "test@yahoo.com")
	user.Add("password", "secret555")
	user.Add("name", "Tester")
	user.Add("facebook_id", "123456")

	response, err := http.PostForm(usersUrl+"/create", user)
	defer response.Body.Close()

	errorIfExists(t, err)
	validateResponseCode(t, http.StatusOK, response.StatusCode)

	responseUser := validateGetUser(t, response.Body, user)

	// update name
	user.Set("name", "No Longer Tester")
	user.Add("userId", fmt.Sprintf("%d", responseUser.user.id))
	updateInvalidResponse, updateInvalidErr := http.PostForm(usersUrl+"/update", user)
	defer updateInvalidResponse.Body.Close()

	// handle error of invalid token
	errorIfExists(t, updateInvalidErr)
	validateResponseCode(t, http.StatusForbidden, updateInvalidResponse.StatusCode)

	// handle successful response
	client := &http.Client{}
	updateRequest, _ := http.NewRequest("POST", usersUrl+"/update", strings.NewReader(user.Encode()))
	updateRequest.Header.Set("Authorization", "Bearer "+responseUser.user.token)
	updateRequest.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	updateResponse, updateErr := client.Do(updateRequest)
	defer updateResponse.Body.Close()

	errorIfExists(t, updateErr)
	validateResponseCode(t, http.StatusOK, updateResponse.StatusCode)
	_ = validateGetUser(t, updateResponse.Body, user)

	// delete user
	// handle successful response
	client = &http.Client{}
	user = url.Values{}
	user.Add("userId", fmt.Sprintf("%d", responseUser.user.id))
	deleteRequest, _ := http.NewRequest("POST", usersUrl+"/delete", strings.NewReader(user.Encode()))
	deleteRequest.Header.Set("Authorization", "Bearer "+responseUser.user.token)
	deleteRequest.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	deleteResponse, deleteErr := client.Do(deleteRequest)
	defer deleteResponse.Body.Close()

	errorIfExists(t, deleteErr)
	validateResponseCode(t, http.StatusOK, deleteResponse.StatusCode)

	// ensure we can't get user or not is_active
	inactiveUserResponse, _ := http.Get(usersUrl + user.Get("userId"))
	defer inactiveUserResponse.Body.Close()
	validateResponseCode(t, http.StatusNotFound, inactiveUserResponse.StatusCode)
}

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

func validateUserAttributeNumeric(t *testing.T, expectedValue int64, value int64, key string) {
	if expectedValue != value {
		t.Errorf("Response user %s changed from %d to %d", key, expectedValue, value)
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

func validateGetUser(t *testing.T, body io.ReadCloser, user url.Values) response {
	var response response
	var responseMap map[string]interface{}
	responseBody, _ := ioutil.ReadAll(body)
	jsonErr := json.Unmarshal(responseBody, &responseMap)

	errorIfExists(t, jsonErr)
	response.code = responseMap["code"].(string)
	responseUserBody := responseMap["response"].(map[string]interface{})
	response.user.id = int64(responseUserBody["id"].(float64))
	response.user.email = responseUserBody["email"].(string)
	response.user.name = responseUserBody["name"].(string)

	if is_active := responseUserBody["is_active"]; is_active != nil {
		response.user.is_active = is_active.(int)
	}

	if created_at := responseUserBody["created_at"]; created_at != nil {
		response.user.created_at = created_at.(string)
	}

	if updated_at := responseUserBody["updated_at"]; updated_at != nil {
		response.user.updated_at = responseUserBody["updated_at"].(string)
	}

	if responseUserBody["token"] != nil {
		response.user.token = responseUserBody["token"].(string)
	}

	validateCode(t, response.code)

	// only validate user id if in request
	validateUserId(t, response.user.id)
	if userId := user.Get("userId"); userId != "" {
		userIdInt, _ := strconv.ParseInt(userId, 10, 64)

		validateUserAttributeNumeric(t, userIdInt, response.user.id, "id")
	}

	if email := user.Get("email"); email != "" {
		validateUserAttribute(t, response.user.email, user.Get("email"), "email")
	}

	if name := user.Get("name"); name != "" {
		validateUserAttribute(t, response.user.name, user.Get("name"), "name")
	}

	// validate doesnt exist
	validateUserPassword(t, responseUserBody["password"])

	return response
}
