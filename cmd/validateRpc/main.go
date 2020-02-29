package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
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
	"strings"
	"time"
)

var (
	localHost        string
	port             = ":8000"
	hashConsistent   *common.Consistent
	checkRightClient pb.CheckRightServiceClient
	accessControl    *common.AccessControl
	rabbitMQValidate *rabbitmq.RabbitMq
	redisClient      *redis.Client
	luaScript        *redis.Script
)

func Check(resp http.ResponseWriter, req *http.Request) {
	fmt.Println("Begin to check.")
	queryForm, err := url.ParseQuery(req.URL.RawQuery)
	if err != nil || len(queryForm) <= 0 {
		resp.Write([]byte("false"))
		return
	}

	productString := queryForm["productID"][0]
	userCookie, err := req.Cookie("uid")
	if err != nil {
		resp.Write([]byte("false"))
		return
	}

	// 根据哈希一致性得到的地址调用checkright服务
	addr, err := GetDistributedAddr(req)
	if err != nil {
		resp.Write([]byte("false"))
		return
	}
	gRpcConn := common.GetGrpcClientConn(addr)
	checkRightClient = pb.NewCheckRightServiceClient(gRpcConn)
	defer gRpcConn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	r, err := checkRightClient.CheckRight(ctx, &pb.Uid{Value: userCookie.Value})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	right := r.Value
	cancel()
	if !right {
		resp.Write([]byte("false"))
		return
	}
	productKey := productString + "_balance"
	res, err := luaScript.Run(redisClient, []string{productKey}, 1).Result()
	if err != nil {
		log.Println("reduce the product balance error:",err)
		resp.Write([]byte("false"))
		return
	}
	balanceStr := res.(string)
	banlanceInt, err := strconv.Atoi(balanceStr)
	if err != nil {
		resp.Write([]byte("false"))
		return
	}
	if banlanceInt > 0 {
		userID, err := strconv.ParseInt(userCookie.Value, 10, 64)
		if err != nil {
			resp.Write([]byte("false"))
			return
		}
		productID, err := strconv.ParseInt(productString, 10, 64)
		message := models.NewMessage(productID, userID)
		byteMessage, err := json.Marshal(message)
		if err != nil {
			resp.Write([]byte("false"))
			return
		}
		err = rabbitMQValidate.PublishSimple(string(byteMessage))
		if err != nil {
			log.Println("publish the msg err：L",err)
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
	uid, err := url.QueryUnescape(uidCookie.Value)
	if err != nil {
		return err
	}
	sign, err := url.QueryUnescape(signCookie.Value)
	if err != nil {
		return err
	}

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

type server struct {
	pb.UnimplementedCheckRightServiceServer
}

//实现CheckOneServiceServer
func (s *server) CheckRight(ctx context.Context, req *pb.Uid) (*pb.IsOk, error) {
	uid := req.GetValue()
	isOk := accessControl.GetDataFromMap(uid)
	return &pb.IsOk{Value: isOk}, nil
}

func GetDistributedAddr(req *http.Request) (string, error) {
	uid, err := req.Cookie("uid")
	if err != nil {
		return "", err
	}
	hostRequest, err := hashConsistent.Get(uid.Value)
	if err != nil {
		return "", err
	}
	return hostRequest + port, nil
}

func main() {
	var err error
	localHost, err = common.GetIntranceIp()
	if err != nil {
		panic(err)
	}
	accessControl = common.GetAccessControl()
	hashConsistent = common.NewConsistent()
	redisClient = common.GetClientFromSen()
	luaScript, err = common.GetDecrbyScr()
	if err != nil {
		log.Panic(err)
	}
	crt, key := common.GetGrpcCrtKey()
	rabbitMQValidate = rabbitmq.NewRabbitMQSimple("miaosha")

	// 创建GetRight的grpc服务
	grpcServer := common.GetGrpcServer()
	pb.RegisterCheckRightServiceServer(grpcServer, &server{})

	for _, v := range configs.Cfg.NODES.Address {
		hashConsistent.Add(v)
	}
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	filter := common.NewFilter()
	filter.RegisterFilterUri("check", Auth)
	mux := http.NewServeMux()
	mux.HandleFunc("/check", filter.Handler(Check))

	http.ListenAndServeTLS(port, crt, key,
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.ProtoMajor != 2 {
				mux.ServeHTTP(w, r)
				return
			}
			if strings.Contains(
				r.Header.Get("Content-Type"), "application/grpc",
			) {
				grpcServer.ServeHTTP(w, r) // gRPC Server
				return
			}
			mux.ServeHTTP(w, r)
			return
		}),
	)
}
