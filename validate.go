package main

import (
	"errors"
	"fmt"
	"net/http"
	"pro-iris/common"
	"pro-iris/encrypt"
)

// cluster addresses : here internal IPs
var hostArray = []string{"127.0.0.1", "127.0.0.1"}

var localHost = "127.0.0.1"

var port = "8081"

var hashConsistent *common.Consistent

// Unified verification filter
// Each interface needs to be verified in advance
func Auth(w http.ResponseWriter, r *http.Request) error {
	fmt.Println("run Auth function successfully")
	// Add permission verification based on cookie
	err := checkUserInfo(r)
	if err != nil {
		return err
	}
	return nil
}

// Identity verification
func checkUserInfo(r *http.Request) error {
	// get uid from cookie
	uidCookie, err := r.Cookie("uid")
	if err != nil {
		return errors.New("unable to access uid from cookie")
	}
	// get encrypted userInfo from cookie
	signCookie, err := r.Cookie("sign")
	if err != nil {
		return errors.New("unable to access sign from cookie")
	}
	// Decrypt sign
	signByte, err := encrypt.DePwdCode(signCookie.Value)
	if err != nil {
		return errors.New("encrypted User information has been tampered with")
	}

	fmt.Println("Identity verification")
	fmt.Println("uid：" + uidCookie.Value)
	fmt.Println("decrypted sign：" + string(signByte))
	if checkInfo(uidCookie.Value, string(signByte)) {
		return nil
	} else {
		return errors.New("identity verification failed")
	}
}

func checkInfo(checkStr string, signStr string) bool {
	return checkStr == signStr
}

// Execute normal logic
func Check(w http.ResponseWriter, r *http.Request) {
	fmt.Println("run Check function successfully")
}

func main() {
	// LB settings
	// Consistent Hashing
	hashConsistent = common.NewConsistent()
	for _, v := range hostArray {
		hashConsistent.Add(v)
	}

	// 1. create filter
	filter := common.NewFilter()
	// 2. register filter
	filter.RegisterFilterUri("/check", Auth)
	// 3. start service
	http.HandleFunc("/check", filter.Handle(Check))
	http.ListenAndServe(":8083", nil)
}
