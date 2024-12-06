# aliyun-oss-appserver-go
Upload data to OSS through Web applications. Add signatures on the server, configure upload callback, and directly transfer data.

## Installation
```shell
go get github.com/alphasnow/aliyun-oss-appserver-go
```

## Usage
### Token generate 
```go
token := appserver.NewToken(&Config{
    AccessKeyId:     "yourAccessKeyId",
    AccessKeySecret: "yourAccessKeySecret",
    Host:            "https://bucket-name.oss-cn-hangzhou.aliyuncs.com",
    UploadDir:       "user-dir-prefix/",
    ExpireTime:      30,
    CallbackUrl:     "http://88.88.88.88:8888",
})
policyToken, _ := token.Generate()
tokenJson, _ := json.Marshal(policyToken)
//{
//    "accessid": "yourAccessKeyId",
//    "host": "https://bucket-name.oss-cn-hangzhou.aliyuncs.com",
//    "expire": 1735689600,
//    "signature": "uXL82wU5IGCd7vcZKX9gua5TUJs=",
//    "policy": "eyJleHBpcmF0aW9uIjoiMjAyNS0wMS0wMVQwMDowMDowMFoiLCJjb25kaXRpb25zIjpbWyJzdGFydHMtd2l0aCIsIiRrZXkiLCJ1c2VyLWRpci1wcmVmaXgvIl1dfQ==",
//    "dir": "user-dir-prefix/",
//    "callback": "eyJjYWxsYmFja1VybCI6Imh0dHA6Ly84OC44OC44OC44ODo4ODg4IiwiY2FsbGJhY2tCb2R5IjoiZmlsZW5hbWU9JHtvYmplY3R9XHUwMDI2c2l6ZT0ke3NpemV9XHUwMDI2bWltZVR5cGU9JHttaW1lVHlwZX1cdTAwMjZoZWlnaHQ9JHtpbWFnZUluZnQuaGVpZ2h0fVx1MDAyNndpZHRoPSR7aW1hZ2VJbmZ0LndpZHRofSIsImNhbGxiYWNrQm9keVR5cGUiOiJhcHBsaWNhdGlvbi94LXd3dy1mb3JtLXVybGVuY29kZWQifQ=="
//}

```

## Callback verify
```go
err := appserver.VerifySignatureByRequest(request)
if err != nil {
    return err
}
bodyContent, _ := io.ReadAll(request.Body)
defer request.Body.Close()
callbackBody := new(appserver.CallbackBody)
json.Unmarshal(bodyContent, callbackBody)
//{
//    "bucket": "bucket-name",
//    "object": "user-dir-prefix/image.jpg",
//    "etag": "A3AC1B2FAADBD*****EE9F5EA57CAACB",
//    "size": 2788,
//    "mimeType": "image/jpeg",
//    "imageInfo": {
//        "height": 197,
//        "width": 257,
//        "format": "jpg"
//    },
//    "crc64": 34616313***72852000,
//    "contentMd5": "o6wbL6rb0***7p9epXyqyw==",
//    "vpcId": null,
//    "clientIp": "222.10.20.30",
//    "reqId": "674EB5AA2*****37341888F8",
//    "operation": "PutObject"
//}
```

## Reference
- reference code 
[aliyun-oss-appserver-go-master.zip](https://help-static-aliyun-doc.aliyuncs.com/file-manage-files/zh-CN/20240710/zbucef/aliyun-oss-appserver-go-master.zip)
- reference doc [https://www.alibabacloud.com/help/en/oss/use-cases/go-1](https://www.alibabacloud.com/help/en/oss/use-cases/go-1)
