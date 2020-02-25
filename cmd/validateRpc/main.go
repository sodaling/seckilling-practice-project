package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"log"
	"net/http"
	"net/url"
	"os"
	"seckilling-practice-project/common"
	"seckilling-practice-project/configs"
	pb "seckilling-practice-project/grpc"
	"seckilling-practice-project/models"
	"seckilling-practice-project/rabbitmq"
	"strconv"
	"sync"
	"time"
)

var localHost string

var port = "8000"

var hashConsistent *common.Consistent

var gRpcAddress = "localhost:50051"

var gRpcClient pb.GetOneServiceClient

type AccessControl struct {
	sourceArray map[int]string
	sync.RWMutex
}

var accessControl AccessControl

var rabbitMQValidate *rabbitmq.RabbitMq

func (m *AccessControl) GetNewRecord(uid int) interface{} {
	m.RLock()
	defer m.RUnlock()
	return m.sourceArray[uid]
}

func (m *AccessControl) GetDistributedRight(req *http.Request) bool {
	uid, err := req.Cookie("uid")
	if err != nil {
		return false
	}
	hostRequest, err := hashConsistent.Get(uid.Value)
	if err != nil {
		return false
	}
	fmt.Println(localHost)
	fmt.Println(hostRequest)
	if hostRequest == localHost {
		return m.GetDataFromMap(uid.Value)
	} else {
		return m.GetDataFromOtherMap(hostRequest, req)
	}
}
func (m *AccessControl) GetDataFromMap(key string) bool {
	//uid, err := strconv.Atoi(key)
	//if err != nil {
	//	return false
	//}
	//if data := m.GetNewRecord(uid); data == nil {
	//	return false
	//} else {
	//	return true
	//}
	return true
}
func (m *AccessControl) SetNewRocord(uid int) {
	m.Lock()
	defer m.Unlock()
	m.sourceArray[uid] = "test"
}

func (m *AccessControl) GetDataFromOtherMap(host string, request *http.Request) bool {
	resp, body, err := common.GetUrl("http://"+host+":"+port+"/checkRight", request)
	if err != nil {
		return false
	}
	if resp.StatusCode == http.StatusOK {
		if string(body) == "true" {
			return true
		} else {
			return false
		}
	}
	return false
}

func Check(resp http.ResponseWriter, req *http.Request) {
	fmt.Println("Begin to check.")
	queryForm, err := url.ParseQuery(req.URL.RawQuery)
	if err != nil || len(queryForm) <= 0 {
		resp.Write([]byte("false"))
		return
	}

	productString := queryForm["productID"][0]
	fmt.Println(productString)
	userCookie, err := req.Cookie("uid")
	if err != nil {
		resp.Write([]byte("false"))
		return
	}

	right := accessControl.GetDistributedRight(req)
	if !right {
		resp.Write([]byte("false"))
		return
	}

	// 通过grpc获得getone结果

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := gRpcClient.GetOne(ctx, &empty.Empty{})
	if err != nil {
		resp.Write([]byte("false"))
		return
	}
	if r.GetValue() {

		productID, err := strconv.ParseInt(productString, 10, 64)
		if err != nil {
			resp.Write([]byte("false"))
			return
		}

		userID, err := strconv.ParseInt(userCookie.Value, 10, 64)
		if err != nil {
			resp.Write([]byte("false"))
			return
		}
		message := models.NewMessage(productID, userID)
		byteMessage, err := json.Marshal(message)
		if err != nil {
			resp.Write([]byte("false"))
			return
		}
		err = rabbitMQValidate.PublishSimple(string(byteMessage))
		if err != nil {
			resp.Write([]byte("false"))
			return
		}
		resp.Write([]byte("true"))
		return
	}
	resp.Write([]byte("false"))
	return
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

	//unescape because of gin cookie operation
	uid, _ := url.QueryUnescape(uidCookie.Value)
	sign, _ := url.QueryUnescape(signCookie.Value)

	signByte, err := common.DePwdCode(sign)
	if err != nil {
		return err
	}
	if checkInfo(uid, string(signByte)) {
		return nil
	}
	return errors.New("auth failed")

}

func checkInfo(checkStr, signStr string) bool {
	return checkStr == signStr
}

func main() {
	hashConsistent = common.NewConsistent()
	for _, v := range configs.Cfg.NODES.Address {
		hashConsistent.Add(v)
	}
	var err error
	localHost, err = common.GetIntranceIp()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// 创建grpc连接
	gRpcConn, err := grpc.Dial(gRpcAddress, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer gRpcConn.Close()
	gRpcClient = pb.NewGetOneServiceClient(gRpcConn)

	rabbitMQValidate = rabbitmq.NewRabbitMQSimple("miaosha")
	filter := common.NewFilter()
	filter.RegisterFilterUri("check", Auth)
	http.HandleFunc("/check", filter.Handler(Check))
	http.HandleFunc("/checkRight", filter.Handler(CheckRight))

	http.ListenAndServe(":8000", nil)
}

func CheckRight(resp http.ResponseWriter, req *http.Request) {
	uidCookie, err := req.Cookie("uid")
	if err != nil {
		resp.Write([]byte("false"))
		return
	}
	uid := uidCookie.Value
	isOk := accessControl.GetDataFromMap(uid)
	if isOk {
		resp.Write([]byte("true"))
	} else {
		resp.Write([]byte("false"))
	}
}
