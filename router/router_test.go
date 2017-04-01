package router

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetRouter(t *testing.T) {
	router := GetRouter()
	routerType := reflect.TypeOf(router)
	handlerType := reflect.TypeOf((*http.Handler)(nil)).Elem()

	assert.True(t, routerType.Implements(handlerType))
}
