package main

import (
	"log"
	"net/http"

	"github.com/pick-up-api/router"
)

func main() {
	log.Fatal(http.ListenAndServe(":3001", router.GetRouter()))
}
