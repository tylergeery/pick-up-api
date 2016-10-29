package response

import (
    "net/http"
    "encoding/json"
)

type Response struct {
    Code string `json:"code"`
    Response interface{} `json:"response"`
}

type Error struct {
    Message string `json:"message"`
}

func Success(w http.ResponseWriter, jsonObject interface{}) {
    response := Response{"ok", jsonObject}

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

func Fail(w http.ResponseWriter, httpStatus int, message string) {
    response := Response{"err", Error{message}}

    w.WriteHeader(httpStatus)
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}
