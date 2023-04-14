package main

import (
	"fmt"
	"net/http"
	"pro-iris/common"
)

// Unified verification filter
// Each interface needs to be verified in advance
func Auth(w http.ResponseWriter, r *http.Request) error {
	fmt.Println("run Auth function successfully")
	// Add permission verification based on cookie
	return nil
}

//func checkUserInfo(r *http.Request)

// Execute normal logic
func Check(w http.ResponseWriter, r *http.Request) {
	fmt.Println("run Check function successfully")
}

func main() {
	// 1. create filter
	filter := common.NewFilter()
	// 2. register filter
	filter.RegisterFilterUri("/check", Auth)
	// 3. start service
	http.HandleFunc("/check", filter.Handle(Check))
	http.ListenAndServe(":8003", nil)
}
