package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"pro-iris/common"
	"pro-iris/encrypt"
	"strconv"
	"sync"
)

// cluster addresses : here internal IPs
var hostArray = []string{"127.0.0.1", "127.0.0.1"}

var localHost = "127.0.0.1"

var port = "8081"

var hashConsistent *common.Consistent

// store control info
type AccessControl struct {
	// Store information based on user ID
	sourcesArray map[int]interface{}
	// Using RW mutex to ensure R/W security of map
	sync.RWMutex
}

var accessControl = &AccessControl{
	sourcesArray: make(map[int]interface{}),
}

func (m *AccessControl) GetNewRecord(uid int) interface{} {
	m.RWMutex.RLock()
	defer m.RWMutex.RUnlock()
	data := m.sourcesArray[uid]
	return data
}

func (m *AccessControl) SetNewRecord(uid int) {
	m.RWMutex.Lock()
	m.sourcesArray[uid] = "hello world"
	m.RWMutex.Unlock()
}

func (m *AccessControl) GetDistributedRight(req *http.Request) bool {
	uid, err := req.Cookie("uid")
	if err != nil {
		return false
	}
	// consistent hashing
	hostRequest, err := hashConsistent.Get(uid.Value)
	if err != nil {
		return false
	}

	if hostRequest == localHost {
		// Local: data reading and verification
		return m.GetDataFromMap(uid.Value)
	} else {
		// using agent
		return GetDataFromOtherMap(hostRequest, req)
	}
}

func (m *AccessControl) GetDataFromMap(uid string) bool {
	uidInt, err := strconv.Atoi(uid)
	if err != nil {
		return false
	}
	data := m.GetNewRecord(uidInt)
	return data != nil
}

func GetDataFromOtherMap(host string, req *http.Request) bool {
	uidPre, err := req.Cookie("uid")
	if err != nil {
		return false
	}
	uidSign, err := req.Cookie("sign")
	if err != nil {
		return false
	}
	// mock interface access
	client := &http.Client{}
	r, err := http.NewRequest("GET", "http://"+host+":"+port+"/access", nil)
	if err != nil {
		return false
	}
	cookieUid := &http.Cookie{
		Name:  "uid",
		Value: uidPre.Value,
		Path:  "/",
	}
	cookieSign := &http.Cookie{
		Name:  "sign",
		Value: uidSign.Value,
		Path:  "/",
	}
	r.AddCookie(cookieUid)
	r.AddCookie(cookieSign)
	response, err := client.Do(r)
	if err != nil {
		return false
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return false
	}
	if response.StatusCode == 200 {
		return string(body) == "true"
	} else {
		return false
	}
}

// Auth : Unified verification filter
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

// Check : Execute normal logic
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
