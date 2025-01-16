package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("######## Welcome to the my Todo App! ########")
	http.HandleFunc("/", helloUser)
	http.ListenAndServe(":5050", nil)
}

func helloUser(w http.ResponseWriter, r *http.Request) {
	greeting := "Hello, User! Welcome to the my Todo App!"
	fmt.Fprintf(w, greeting)
}
