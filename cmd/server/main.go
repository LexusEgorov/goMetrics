package main

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Message string `json:"message"`
}

func testHandle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := Response{Message: "Hello, Go!"}

	responseJSON, err := json.Marshal(response)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Write(responseJSON)
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc(`/`, testHandle)

	http.ListenAndServe(":8080", mux)
}
