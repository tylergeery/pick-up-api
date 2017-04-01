package controllerTests

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"

	"github.com/gorilla/mux"
	"github.com/pick-up-api/router"
	"github.com/pick-up-api/utils/messaging"
	"github.com/pick-up-api/utils/resources"
	"github.com/pick-up-api/utils/validation"
)

type user struct {
	id            int64
	email         string
	name          string
	refresh_token string
	access_token  string
	is_active     int
	created_at    string
	updated_at    string
}

type response struct {
	code string
	user user
}

var (
	db           *sql.DB
	handler      *mux.Router
	userEndpoint string
	authEndpoint string
)

func init() {
	db = resources.DB()
	handler = router.GetRouter()
	userEndpoint = "https://api.pickup.com/user"
	authEndpoint = "https://api.pickup.com/auth"
}

func setUp(t *testing.T) {
	resources.SetTestTXInterface()
}

func tearDown(t *testing.T) {
	_ = resources.TX().TearDown()
}

func TestUserProfile(t *testing.T) {
	setUp(t)

	// Given
	validUserRecorder := httptest.NewRecorder()
	validReq, validErr := http.NewRequest("GET", userEndpoint+"/1", nil)
	outOfBoundsRecorder := httptest.NewRecorder()
	outOfBoundsReq, outOfBoundsErr := http.NewRequest("GET", userEndpoint+"/11234123", nil)

	// Perform
	errorIfExists(t, validErr)
	errorIfExists(t, outOfBoundsErr)

	handler.ServeHTTP(validUserRecorder, validReq)
	handler.ServeHTTP(outOfBoundsRecorder, outOfBoundsReq)

	// Assert
	validateResponseCode(t, http.StatusOK, validUserRecorder.Result().StatusCode)
	validateResponseCode(t, http.StatusNotFound, outOfBoundsRecorder.Result().StatusCode)

	_ = validateGetUser(t, validUserRecorder.Body.Bytes(), url.Values{})

	tearDown(t)
}

func TestCreateUser(t *testing.T) {
	setUp(t)

	// Given
	emptyUser := url.Values{}

	userWithoutEmail := url.Values{}
	userWithoutEmail.Add("password", "secret")
	userWithoutEmail.Add("name", "Tester")

	userWithInvalidEmail := url.Values{}
	userWithInvalidEmail.Add("email", "test.test")
	userWithInvalidEmail.Add("password", "secret555")
	userWithInvalidEmail.Add("name", "Tester")

	userWithInvalidPassword := url.Values{}
	userWithInvalidPassword.Add("email", "test.test@test.com")
	userWithInvalidPassword.Add("password", "secret")
	userWithInvalidPassword.Add("name", "Tester")

	userWithoutPassword := url.Values{}
	userWithoutPassword.Add("email", "test.test")
	userWithoutPassword.Add("name", "Tester")

	user := url.Values{}
	user.Add("email", "test+unique55@yahoo.com")
	user.Add("password", "secret555")
	user.Add("name", "Tester")

	fails := []url.Values{
		emptyUser,
		userWithoutEmail,
		userWithInvalidEmail,
		userWithInvalidPassword,
		userWithoutPassword}
	failMessages := []string{
		messaging.USER_REQUIRES_EMAIL,
		messaging.USER_REQUIRES_EMAIL,
		messaging.USER_EMAIL_INVALID,
		validation.PasswordShort,
		messaging.USER_REQUIRES_PASSWORD}
	successes := []url.Values{user}

	// Perform && Assert
	// test invalid params
	for i, invalidUser := range fails {
		invalidCreateRecorder := httptest.NewRecorder()
		invalidReq, invalidErr := http.NewRequest(
			"POST",
			userEndpoint+"/create",
			bytes.NewBufferString(invalidUser.Encode()))

		errorIfExists(t, invalidErr)
		invalidReq.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		handler.ServeHTTP(invalidCreateRecorder, invalidReq)
		validateResponseCode(t, http.StatusBadRequest, invalidCreateRecorder.Result().StatusCode)

		invalidCreateResponse := getResponse(t, invalidCreateRecorder.Body.Bytes())["response"].(map[string]interface{})
		validateResponseMessage(t, failMessages[i], invalidCreateResponse["message"].(string))
		log.Printf("Fail message: %s\n", invalidCreateResponse["message"].(string))
	}

	// test valid user json
	for _, validUser := range successes {
		validCreateRecorder := httptest.NewRecorder()
		validReq, validErr := http.NewRequest(
			"POST",
			userEndpoint+"/create",
			bytes.NewBufferString(validUser.Encode()))

		errorIfExists(t, validErr)
		validReq.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		handler.ServeHTTP(validCreateRecorder, validReq)
		validateResponseCode(t, http.StatusOK, validCreateRecorder.Result().StatusCode)

		_ = validateGetUser(t, validCreateRecorder.Body.Bytes(), validUser)
	}

	tearDown(t)
}

func TestCreateUpdateAndDeleteUser(t *testing.T) {
	setUp(t)

	user := url.Values{}
	user.Add("email", "test+unique22@yahoo.com")
	user.Add("password", "secret555")
	user.Add("name", "Tester")
	user.Add("facebook_id", "123456")

	createRecorder := httptest.NewRecorder()
	createReq, createErr := http.NewRequest(
		"POST",
		userEndpoint+"/create",
		bytes.NewBufferString(user.Encode()))

	errorIfExists(t, createErr)
	createReq.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	handler.ServeHTTP(createRecorder, createReq)
	validateResponseCode(t, http.StatusOK, createRecorder.Result().StatusCode)

	responseUser := validateGetUser(t, createRecorder.Body.Bytes(), user)

	// retrieve valid access token from auth endpoint
	authRecorder := httptest.NewRecorder()
	authReq, authErr := http.NewRequest(
		"POST",
		authEndpoint+"/refresh",
		bytes.NewBufferString(""),
	)

	// get access token
	errorIfExists(t, authErr)
	authReq.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	authReq.Header.Set("Authorization", "Bearer "+responseUser.user.refresh_token)
	handler.ServeHTTP(authRecorder, authReq)
	validateResponseCode(t, http.StatusOK, authRecorder.Result().StatusCode)

	authUser := validateGetUser(t, authRecorder.Body.Bytes(), user)

	// update name
	user.Set("name", "No Longer Tester")
	user.Add("userId", fmt.Sprintf("%d", responseUser.user.id))
	invalidUpdateRecorder := httptest.NewRecorder()
	invalidUpdateReq, invalidUpdateErr := http.NewRequest(
		"POST",
		userEndpoint+"/update",
		bytes.NewBufferString(user.Encode()))

	// handle error of invalid token
	errorIfExists(t, invalidUpdateErr)
	invalidUpdateReq.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	handler.ServeHTTP(invalidUpdateRecorder, invalidUpdateReq)
	validateResponseCode(t, http.StatusForbidden, invalidUpdateRecorder.Result().StatusCode)

	invalidTokenResponse := getResponse(t, invalidUpdateRecorder.Body.Bytes())["response"].(map[string]interface{})
	validateResponseMessage(t, messaging.USER_UNAUTHORIZED_UPDATE, invalidTokenResponse["message"].(string))

	// handle successful response
	updateRecorder := httptest.NewRecorder()
	updateReq, updateErr := http.NewRequest(
		"POST",
		userEndpoint+"/update",
		bytes.NewBufferString(user.Encode()))

	// handle successful update
	errorIfExists(t, updateErr)
	updateReq.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	updateReq.Header.Set("Authorization", "Bearer "+authUser.user.access_token)
	handler.ServeHTTP(updateRecorder, updateReq)
	validateResponseCode(t, http.StatusOK, updateRecorder.Result().StatusCode)

	_ = validateGetUser(t, updateRecorder.Body.Bytes(), user)

	// delete user
	// handle successful response
	user = url.Values{}
	user.Add("userId", fmt.Sprintf("%d", responseUser.user.id))
	deleteRecorder := httptest.NewRecorder()
	deleteReq, deleteErr := http.NewRequest(
		"POST",
		userEndpoint+"/delete",
		bytes.NewBufferString(user.Encode()))

	errorIfExists(t, deleteErr)
	deleteReq.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	deleteReq.Header.Set("Authorization", "Bearer "+authUser.user.access_token)
	handler.ServeHTTP(deleteRecorder, deleteReq)
	validateResponseCode(t, http.StatusOK, deleteRecorder.Result().StatusCode)

	// ensure we can't get user or not is_active
	userRecorder := httptest.NewRecorder()
	userReq, userErr := http.NewRequest(
		"GET",
		fmt.Sprintf("%s/%d", userEndpoint, responseUser.user.id),
		nil)

	errorIfExists(t, userErr)
	handler.ServeHTTP(userRecorder, userReq)

	validateResponseCode(t, http.StatusNotFound, userRecorder.Result().StatusCode)

	notActiveResponse := getResponse(t, userRecorder.Body.Bytes())["response"].(map[string]interface{})
	validateResponseMessage(t, messaging.USER_NOT_ACTIVE, notActiveResponse["message"].(string))

	tearDown(t)
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

func validateResponseMessage(t *testing.T, expectedMessage string, responseMessage string) {
	if expectedMessage != responseMessage {
		t.Errorf("Expected Response message (%s), received (%s)", expectedMessage, responseMessage)
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

func getResponse(t *testing.T, responseBody []byte) map[string]interface{} {
	var responseMap map[string]interface{}
	jsonErr := json.Unmarshal(responseBody, &responseMap)

	errorIfExists(t, jsonErr)

	return responseMap
}

func validateGetUser(t *testing.T, responseBody []byte, user url.Values) response {
	var response response
	responseMap := getResponse(t, responseBody)

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

	if responseUserBody["refresh_token"] != nil {
		response.user.refresh_token = responseUserBody["refresh_token"].(string)
	}

	if responseUserBody["access_token"] != nil {
		response.user.access_token = responseUserBody["access_token"].(string)
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
