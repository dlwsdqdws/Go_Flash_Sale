package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

var sum int64 = 0

// stored product quantity
var productNum int64 = 1000000

var mutex sync.Mutex

// rabbitmq queue request rate
var count int64 = 0

// GetOneProduct get flash sale product
func GetOneProduct() bool {
	mutex.Lock()
	defer mutex.Unlock()
	count += 1
	// Set to release a product every 100 requests
	//if count%100 == 0 {  // delete this line when local test
	if sum < productNum {
		sum += 1
		fmt.Println(sum)
		return true
	}
	//}
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
	err := http.ListenAndServe(":8084", nil)
	if err != nil {
		log.Fatal("Error Occurred:", err)
	}
}
