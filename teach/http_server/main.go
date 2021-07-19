package main

import (
	"log"
	"net/http"
)



func main() {

	http.HandleFunc("/self", func(writer http.ResponseWriter, request *http.Request) {
		log.Println("test")
		writer.Write([]byte("Hello, world"))
	})

	http.ListenAndServe(":8082",nil)
}