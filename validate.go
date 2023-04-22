package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"pro-iris/common"
	"pro-iris/datamodels"
	"pro-iris/encrypt"
	"pro-iris/rabbitmq"
	"strconv"
	"sync"
	"time"
)

// cluster addresses : here internal IPs
//var hostArray = []string{"172.26.194.41", "172.26.194.42"}  // for real server
var hostArray = []string{"192.168.68.105", "192.168.68.105"} // for local test

var localHost = ""

// GetOneIp : getOne SLB intranet IP
//var GetOneIp = "172.26.194.43"  // for real server
var GetOneIp = "127.0.0.1" // for local test

var port = "8083"

var GetOnePort = "8084"

var hashConsistent *common.Consistent

var rabbitMqValidate *rabbitmq.RabbitMQ

// Set server response interval to be 10 seconds -> Avoid malicious for-loop requests
var interval = 10

type BlackList struct {
	listArray map[int]bool
	sync.RWMutex
}

var blackList = &BlackList{
	listArray: make(map[int]bool),
}

func (m *BlackList) GetBlackListByID(uid int) bool {
	m.RLock()
	defer m.RUnlock()
	return m.listArray[uid]
}

func (m *BlackList) SetBlackListByID(uid int) bool {
	m.Lock()
	defer m.Unlock()
	m.listArray[uid] = true
	return true
}

// AccessControl : store control info
type AccessControl struct {
	// Store information based on user ID
	sourcesArray map[int]time.Time
	// Using RW mutex to ensure R/W security of map
	sync.RWMutex
}

var accessControl = &AccessControl{
	sourcesArray: make(map[int]time.Time),
}

func (m *AccessControl) GetNewRecord(uid int) time.Time {
	m.RWMutex.RLock()
	defer m.RWMutex.RUnlock()
	data := m.sourcesArray[uid]
	return data
}

func (m *AccessControl) SetNewRecord(uid int) {
	m.RWMutex.Lock()
	m.sourcesArray[uid] = time.Now()
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
	if blackList.GetBlackListByID(uidInt) {
		return false
	}
	data := m.GetNewRecord(uidInt)
	if !data.IsZero() {
		if data.Add(time.Duration(interval) * time.Second).After(time.Now()) {
			return false
		}
	}
	m.SetNewRecord(uidInt)
	return true
}

func GetDataFromOtherMap(host string, req *http.Request) bool {
	hostUrl := "http://" + host + ":" + port + "/checkRight"
	response, body, err := GetCurl(hostUrl, req)
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
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8082")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
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

func CheckRight(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8082")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	right := accessControl.GetDistributedRight(r)
	if !right {
		w.Write([]byte("false"))
		return
	}
	w.Write([]byte("true"))
	return
}

// Check : Execute normal logic
func Check(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8082")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	fmt.Println("run Check function successfully")
	queryForm, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil || len(queryForm["productID"]) <= 0 {
		w.Write([]byte("false"))
		return
	}
	productString := queryForm["productID"][0]
	fmt.Println(productString)
	userCookie, err := r.Cookie("uid")
	if err != nil {
		w.Write([]byte("false"))
		return
	}

	// 1.Distributed permission verification
	right := accessControl.GetDistributedRight(r)
	if right == false {
		w.Write([]byte("false"))
		return
	}
	// 2.getOne quantity control : prevent oversold
	hostUrl := "http://" + GetOneIp + ":" + GetOnePort + "/getOne"
	responseValidate, bodyValidate, err := GetCurl(hostUrl, r)
	if err != nil {
		w.Write([]byte("false"))
		return
	}
	// 3.check request status of quantity control
	if responseValidate.StatusCode == 200 {
		if string(bodyValidate) == "true" {
			productID, err := strconv.ParseInt(productString, 10, 64)
			if err != nil {
				w.Write([]byte("false"))
				return
			}
			userID, err := strconv.ParseInt(userCookie.Value, 10, 64)
			if err != nil {
				w.Write([]byte("false"))
				return
			}
			message := datamodels.NewMessage(userID, productID)
			byteMsg, err := json.Marshal(message)
			if err != nil {
				w.Write([]byte("false"))
				return
			}
			err = rabbitMqValidate.PublishSimple(string(byteMsg))
			if err != nil {
				w.Write([]byte("false"))
				return
			}
			w.Write([]byte("true"))
			// Every user can buy only once
			//blackList.SetBlackListByID(userID)
			return
		}
	}
	w.Write([]byte("false"))
	return
}

func GetCurl(hostUrl string, request *http.Request) (response *http.Response, body []byte, err error) {
	uidPre, err := request.Cookie("uid")
	if err != nil {
		return
	}
	uidSign, err := request.Cookie("sign")
	if err != nil {
		return
	}
	client := &http.Client{}
	req, err := http.NewRequest("GET", hostUrl, nil)
	if err != nil {
		return
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
	req.AddCookie(cookieUid)
	req.AddCookie(cookieSign)
	response, err = client.Do(req)
	defer response.Body.Close()
	if err != nil {
		return
	}
	body, err = ioutil.ReadAll(response.Body)
	return
}

func main() {
	// LB settings
	// Consistent Hashing
	hashConsistent = common.NewConsistent()
	for _, v := range hostArray {
		hashConsistent.Add(v)
	}

	localIp, err := common.GetIntranetIp()
	if err != nil {
		fmt.Println(err)
	}
	localHost = localIp
	fmt.Println(localHost)

	rabbitMqValidate = rabbitmq.NewRabbitMQSimple("rabbitmqProduct")
	defer rabbitMqValidate.Destroy()

	// 1. create filter
	filter := common.NewFilter()
	// 2. register filter
	filter.RegisterFilterUri("/check", Auth)
	filter.RegisterFilterUri("/checkRight", Auth)
	// 3. start service
	http.HandleFunc("/check", filter.Handle(Check))
	http.HandleFunc("/checkRight", filter.Handle(CheckRight))
	http.ListenAndServe(":8083", nil)
}
