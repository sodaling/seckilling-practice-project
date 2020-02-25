# 保存定义gRPC服务的proto文件

根据proto文件生成go文件：

```shell
$ protoc --go_out=plugins=grpc:. getone.proto
$ protoc --go_out=plugins=grpc:. checkright.proto
```

