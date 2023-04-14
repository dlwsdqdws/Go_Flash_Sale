package main

import (
	"fmt"
	"net/http"
	"pro-iris/common"
)

// Unified verification filter
// Each interface needs to be verified in advance
func Auth(w http.ResponseWriter, r *http.Response) error {
	fmt.Println("run Auth function successfully")
	return nil
}

// Execute normal logic
func Check(w http.ResponseWriter, r *http.Response) {
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
