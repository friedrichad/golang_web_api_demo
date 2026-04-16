package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/demo", demoHandler)
	log.Println("Server starting from log...")
	fmt.Println("Server from fmt...")
	error := http.ListenAndServe(":8080", nil)
	if error != nil {
		log.Fatal("Error starting server:", error)
	}
}
func demoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	response := map[string]string{
		"message": "Hello, this is a json response!",
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("x-Course", "Golang")

	data, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Write(data)

}
