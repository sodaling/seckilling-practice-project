# 生成gRPC相关证书

相关后缀名注释:

1. 证书(Certificate)：**.cer**(windows), **.crt**
2. 私钥(Private Key)：**.key**
3. 证书签名请求(Certificate sign request)：**.csr**

#### 先生成CA证书

```shell
$ openssl genrsa -out ca.key 2048
$ openssl req -new -x509 -days 3650 \
    -subj "/C=GB/L=China/O=sodaling/CN=github.com" \
    -key ca.key -out ca.crt
```

#### 然后生成服务端证书，其中利用刚刚的ca.key证书签名

```shell
$ openssl genrsa -out server.key 2048
$ openssl req -new \
    -subj "/C=GB/L=China/O=server/CN=server.io" \
    -key server.key \
    -out server.csr
$ openssl x509 -req -sha256 \
    -CA ca.crt -CAkey ca.key -CAcreateserial -days 3650 \
    -in server.csr \
    -out server.crt
```

#### 最后生产client证书,也利用到前面的ca证书进行签名

```shell
$ openssl genrsa -out client.key 2048
$ openssl req -new \
    -subj "/C=GB/L=China/O=client/CN=client.io" \
    -key client.key \
    -out client.csr
$ openssl x509 -req -sha256 \
    -CA ca.crt -CAkey ca.key -CAcreateserial -days 3650 \
    -in client.csr \
    -out client.crt
```

