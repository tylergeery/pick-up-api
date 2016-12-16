package controllerTests

import (
    "fmt"
    "testing"
    "net/http/httptest"
    "github.com/pick-up-api/controllers"
)

func TestUserProfile(t *testing.T) {
    req := httptest.NewRequest("GET", "http://pickup.com/user/1", nil)
	w := httptest.NewRecorder()
	controllers.UserProfile(w, req)

	fmt.Printf("%d - %s", w.Code, w.Body.String())
}

func TestCreateUser(t *testing.T) {
    // test invalid params
}
