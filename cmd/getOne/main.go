package main

import (
	"fmt"
	"net/http"
	"sync"
)

var ProductSum int64 = 10000
var Sum int64 = 0
var mutex sync.Mutex

func main() {
	http.HandleFunc("/getOne", GetProduct)
	http.ListenAndServe(":8084", nil)
}

func GetProduct(writer http.ResponseWriter, request *http.Request) {
	if GetOneProduct() {
		writer.Write([]byte("true"))
	} else {
		writer.Write([]byte("false"))
	}
	fmt.Println(Sum)
}

func GetOneProduct() bool {
	mutex.Lock()
	defer mutex.Unlock()
	if Sum < ProductSum {
		Sum += 1
		return true
	}
	return false
}
