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
)

// cluster addresses : here internal IPs
var hostArray = []string{"172.26.194.41", "172.26.194.42"}

var localHost = ""

// GetOneIp : getOne SLB intranet IP
var GetOneIp = "172.26.194.43"

var port = "8083"

var GetOnePort = "8084"

var hashConsistent *common.Consistent

var rabbitMqValidate *rabbitmq.RabbitMQ

// AccessControl : store control info
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
	//uidInt, err := strconv.Atoi(uid)
	//if err != nil {
	//	return false
	//}
	//data := m.GetNewRecord(uidInt)
	//return data != nil
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
