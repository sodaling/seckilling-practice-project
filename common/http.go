package common

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func GetUrl(url string, req *http.Request) (*http.Response, []byte, error) {
	uidCookie, err := req.Cookie("uid")
	if err != nil {
		return &http.Response{}, nil, err
	}
	signCookie, err := req.Cookie("sign")
	if err != nil {
		return &http.Response{}, nil, err
	}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return &http.Response{}, nil, err
	}

	cookieUid := &http.Cookie{Name: "uid", Value: uidCookie.Value, Path: "/"}
	cookieSign := &http.Cookie{Name: "sign", Value: signCookie.Value, Path: "/"}
	//添加cookie到模拟的请求中
	request.AddCookie(cookieUid)
	request.AddCookie(cookieSign)

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		fmt.Println(err)
		return &http.Response{}, nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	return resp, body, err
}
