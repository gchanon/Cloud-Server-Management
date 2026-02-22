package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	port := "3000"
	fs := http.FileServer(http.Dir("."))
	http.Handle("/", fs)

	fmt.Printf("Frontend server running on http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
