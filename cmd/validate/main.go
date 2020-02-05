package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"seckilling-practice-project/common"
	"strconv"
	"sync"
)

var hostArray = []string{"127.0.0.1", "127.0.0.1"}

var localHost = "127.0.0.1"

var port = "8081"

var hashConsistent *common.Consistent

type AccessControl struct {
	sourceArray map[int]string
	sync.RWMutex
}

var accessControl AccessControl

func (m *AccessControl) GetNewRecord(uid int) interface{} {
	m.RLock()
	defer m.RUnlock()
	return m.sourceArray[uid]
}

func (m *AccessControl) GetDistuibutedRight(req *http.Request) bool {
	uid, err := req.Cookie("uid")
	if err != nil {
		return false
	}
	hostRequest, err := hashConsistent.Get(uid.Value)
	if err != nil {
		return false
	}

	if hostRequest == localHost {
		return m.GetDataFromMap(uid.Value)
	} else {
		return m.GetDataFromOtherMap(hostRequest, req)
	}
}
func (m *AccessControl) GetDataFromMap(key string) bool {
	uid, err := strconv.Atoi(key)
	if err != nil {
		return false
	}
	if data := m.GetNewRecord(uid); data == nil {
		return false
	} else {
		return true
	}
}
func (m *AccessControl) SetNewRocord(uid int) {
	m.Lock()
	defer m.Unlock()
	m.sourceArray[uid] = "test"
}

func (m *AccessControl) GetDataFromOtherMap(host string, request *http.Request) bool {
	uidCookie, err := request.Cookie("uid")
	if err != nil {
		return false
	}
	signCookie, err := request.Cookie("sign")
	if err != nil {
		return false
	}

	req, err := http.NewRequest("GET", "http://"+host+":"+port+"/check", nil)
	if err != nil {
		return false
	}

	cookieUid := &http.Cookie{Name: "uid", Value: uidCookie.Value, Path: "/"}
	cookieSign := &http.Cookie{Name: "sign", Value: signCookie.Value, Path: "/"}
	//添加cookie到模拟的请求中
	req.AddCookie(cookieUid)
	req.AddCookie(cookieSign)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false
	}

	if resp.StatusCode == http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return false
		}

		if string(body) == "true" {
			return true
		} else {
			return false
		}
	}
	return false

}

func Check(resp http.ResponseWriter, req *http.Request) {
	resp.Write([]byte("okokokok"))
}

func Auth(resp http.ResponseWriter, req *http.Request) error {
	fmt.Println("Auth begin")
	return CheckUserInfo(req)
}

func CheckUserInfo(req *http.Request) error {
	uidCookie, err := req.Cookie("uid")
	if err != nil {
		return err
	}
	signCookie, err := req.Cookie("sign")
	if err != nil {
		return err
	}
	signByte, err := common.DePwdCode(signCookie.Value)
	if err != nil {
		return err
	}
	if checkInfo(uidCookie.Value, string(signByte)) {
		return nil
	}
	return errors.New("auth failed")

}

func checkInfo(checkStr, signStr string) bool {
	return checkStr == signStr
}

func main() {
	hashConsistent = common.NewConsistent()
	for _, v := range hostArray {
		hashConsistent.Add(v)
	}

	filter := common.NewFilter()
	filter.RegisterFilterUri("check", Auth)
	http.HandleFunc("/check", filter.Handler(Check))

	http.ListenAndServe(":8000", nil)
}
