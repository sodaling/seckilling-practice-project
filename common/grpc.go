package common

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	"log"
	"path"
	"seckilling-practice-project/configs"
)

func GetServerCreds() (credentials.TransportCredentials, error) {
	projectPath := configs.GetProjectPath()
	credsPath := path.Join(projectPath, "keys", "server")
	certificate, err := tls.LoadX509KeyPair(path.Join(credsPath, "server.crt"), path.Join(credsPath, "server.key"))
	grpcHandleErr("get sever creds errors:", err)
	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile(path.Join(projectPath, "keys", "ca", "ca.crt"))
	if err != nil {
		log.Fatal(err)
	}
	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		log.Fatal("failed to append certs")
	}

	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{certificate},
		ClientAuth:   tls.RequireAndVerifyClientCert, // NOTE: this is optional!
		ClientCAs:    certPool,
	})
	return creds, nil
}

func GetClientCreds() (credentials.TransportCredentials, error) {
	projectPath := configs.GetProjectPath()
	credsPath := path.Join(projectPath, "keys", "client")
	certificate, err := tls.LoadX509KeyPair(path.Join(credsPath, "client.crt"), path.Join(credsPath, "client.key"))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile(path.Join(projectPath, "keys", "ca", "ca.crt"))
	if err != nil {
		return nil, err
	}
	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		err := errors.New("failed to append ca certs")
		log.Fatalln(err)
		return nil, err
	}
	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{certificate},
		ServerName:   "server.io", // NOTE: this is required!
		RootCAs:      certPool,
	})
	return creds, nil
}

func GetGrpcServer() *grpc.Server {
	// 初始化grpc服务端连接
	creds, err := GetServerCreds()
	grpcHandleErr("init grpc client error:", err)
	s := grpc.NewServer(grpc.Creds(creds))
	return s
}

func GetGrpcClientConn(address string) *grpc.ClientConn {
	// 初始化grpc的client连接
	creds, err := GetClientCreds()
	grpcHandleErr("init grpc client error:", err)
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(creds))
	grpcHandleErr("did not connect:", err)
	return conn
}

func grpcHandleErr(str string, err error) {
	// 在初始化grpc遇到error都作为panic处理
	if err != nil {
		panic(str + err.Error())
	}
}
