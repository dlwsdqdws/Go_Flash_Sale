package main

import (
	"log"
	"net/http"
	"sync"
)

var sum int64 = 0

// stored product quantity
var productNum int64 = 10000

var mutex sync.Mutex

// GetOneProduct get flash sale product
func GetOneProduct() bool {
	mutex.Lock()
	defer mutex.Unlock()
	if sum < productNum {
		sum += 1
		return true
	}
	return false
}

func GetProduct(w http.ResponseWriter, req *http.Request) {
	if GetOneProduct() {
		w.Write([]byte("true"))
		return
	}
	w.Write([]byte("false"))
	return
}

func main() {
	http.HandleFunc("/getOne", GetProduct)
	err := http.ListenAndServe(":12345", nil)
	if err != nil {
		log.Fatal("Error Occurred:", err)
	}
}
