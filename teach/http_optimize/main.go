package main

import (
	"log"
	"main.go/teach/http_optimize/pool"
	"net/http"
)

func main() {
	newPool, e := pool.NewPool(5000000)
	if e != nil {
		log.Println(e)
	}
	defer newPool.Close()

	http.HandleFunc("/optimize", func(writer http.ResponseWriter, request *http.Request) {
		newPool.Put(&pool.Task{
			Handle: func(v ...interface{}) {
				log.Println("do something")
			},
		})
		writer.Write([]byte("success"))
	})

	http.ListenAndServe(":8083", nil)
}
