package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("HEEEEELLLLLOOOOOOOOO Itinitary Lions🦁🦁🦁")

	port := ":3000"
	mux := http.NewServeMux()

	// All Endpoints and Handler calls

	fmt.Println("Starting the server...")
	fmt.Printf("Server running at https://localhost%s\n", port)

	http.ListenAndServe(port, nil)

}
