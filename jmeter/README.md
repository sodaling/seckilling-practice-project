# Jmeter

`cmd/gen-test-data/main.go`是往数据库简单mock了1000个测试用户和10个商品的脚本。同时值得注意的是，由于所有页面（包括分布式互相通信的grpc）也是利用https的，所以jmeter需要另外配置http的请求。相关参考文档如下：

https://www.jianshu.com/p/efe03b4d02a3

https://www.jianshu.com/p/0e4daecc8122