English | [简体中文](README-CN.md)

# aliyun-oss-appserver-go

[![Ask DeepWiki](https://deepwiki.com/badge.svg)](https://deepwiki.com/alphasnow/aliyun-oss-appserver-go)

Upload data to OSS through Web applications. Add signatures on the server, configure upload callback, and directly transfer data.

## Installation

```shell
go get -u github.com/alphasnow/aliyun-oss-appserver-go
```

## Usage

### Token generate

```go
token := appserver.NewToken(&appserver.Config{
    AccessKeyId:     "yourAccessKeyId",
    AccessKeySecret: "yourAccessKeySecret",
    Host:            "https://bucket-name.oss-cn-hangzhou.aliyuncs.com",
    Directory:       "user-dir-prefix/",
    ExpireSecond:    600,
    // Only the PutObject, PostObject, and CompleteMultipartUpload support Callback
    CallbackUrl:     "http://domain.com/oss/callback",
})
postToken, _ := token.Generate()
postTokenJson, _ := json.Marshal(postToken)
//{
//    "OSSAccessKeyId": "yourAccessKeyId",
//    "policy": "eyJleHBpcmF0aW9uIjoiMjAyNS0wMS0wMVQwMDowMDowMFoiLCJjb25kaXRpb25zIjpbWyJzdGFydHMtd2l0aCIsIiRrZXkiLCJ1c2VyLWRpci1wcmVmaXgvIl1dfQ==",
//    "callback": "eyJjYWxsYmFja1VybCI6Imh0dHA6Ly9kb21haW4uY29tL29zcy9jYWxsYmFjayIsImNhbGxiYWNrQm9keSI6IntcImJ1Y2tldFwiOiR7YnVja2V0fSxcIm9iamVjdFwiOiR7b2JqZWN0fSxcImV0YWdcIjoke2V0YWd9LFwic2l6ZVwiOiR7c2l6ZX0sXCJtaW1lVHlwZVwiOiR7bWltZVR5cGV9LFwiaW1hZ2VJbmZvXCI6e1wiaGVpZ2h0XCI6JHtpbWFnZUluZm8uaGVpZ2h0fSxcIndpZHRoXCI6JHtpbWFnZUluZm8ud2lkdGh9LFwiZm9ybWF0XCI6JHtpbWFnZUluZm8uZm9ybWF0fX0sXCJjcmM2NFwiOiR7Y3JjNjR9LFwiY29udGVudE1kNVwiOiR7Y29udGVudE1kNX0sXCJ2cGNJZFwiOiR7dnBjSWR9LFwiY2xpZW50SXBcIjoke2NsaWVudElwfSxcInJlcUlkXCI6JHtyZXFJZH0sXCJvcGVyYXRpb25cIjoke29wZXJhdGlvbn19IiwiY2FsbGJhY2tCb2R5VHlwZSI6ImFwcGxpY2F0aW9uL2pzb24ifQ==",
//    "signature": "uXL82wU5IGCd7vcZKX9gua5TUJs=",
//    "host": "https://bucket-name.oss-cn-hangzhou.aliyuncs.com",
//    "expire": 1735689600,
//    "directory": "user-dir-prefix/"
//}
```

## Upload file

```bash
curl --location "https://bucket-name.oss-cn-hangzhou.aliyuncs.com" \
--form 'key="user-dir-prefix/${filename}"' \
--form 'policy="eyJleHBpcmF0aW9uIjoiMjAyNS0wMS0wMVQwMDowMDowMFoiLCJjb25kaXRpb25zIjpbWyJzdGFydHMtd2l0aCIsIiRrZXkiLCJ1c2VyLWRpci1wcmVmaXgvIl1dfQ=="' \
--form 'OSSAccessKeyId="yourAccessKeyId"' \
--form 'callback="eyJjYWxsYmFja1VybCI6Imh0dHA6Ly9kb21haW4uY29tL29zcy9jYWxsYmFjayIsImNhbGxiYWNrQm9keSI6IntcImJ1Y2tldFwiOiR7YnVja2V0fSxcIm9iamVjdFwiOiR7b2JqZWN0fSxcImV0YWdcIjoke2V0YWd9LFwic2l6ZVwiOiR7c2l6ZX0sXCJtaW1lVHlwZVwiOiR7bWltZVR5cGV9LFwiaW1hZ2VJbmZvXCI6e1wiaGVpZ2h0XCI6JHtpbWFnZUluZm8uaGVpZ2h0fSxcIndpZHRoXCI6JHtpbWFnZUluZm8ud2lkdGh9LFwiZm9ybWF0XCI6JHtpbWFnZUluZm8uZm9ybWF0fX0sXCJjcmM2NFwiOiR7Y3JjNjR9LFwiY29udGVudE1kNVwiOiR7Y29udGVudE1kNX0sXCJ2cGNJZFwiOiR7dnBjSWR9LFwiY2xpZW50SXBcIjoke2NsaWVudElwfSxcInJlcUlkXCI6JHtyZXFJZH0sXCJvcGVyYXRpb25cIjoke29wZXJhdGlvbn19IiwiY2FsbGJhY2tCb2R5VHlwZSI6ImFwcGxpY2F0aW9uL2pzb24ifQ=="' \
--form 'signature="uXL82wU5IGCd7vcZKX9gua5TUJs="' \
--form 'file=@"~/Downloads/image.jpg"'
```

## Callback verify

```go
aliyunOSSCallback := appserver.NewAliyunOSSCallback(request)
callbackBody,err := aliyunOSSCallback.VerifySignature()
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
//    "clientIp": "100.20.30.40",
//    "reqId": "674EB5AA2*****37341888F8",
//    "operation": "PutObject"
//}
```

## Reference

- reference code [aliyun-oss-appserver-go-master.zip](https://help-static-aliyun-doc.aliyuncs.com/file-manage-files/zh-CN/20240710/zbucef/aliyun-oss-appserver-go-master.zip)
- reference doc [https://www.alibabacloud.com/help/en/oss/use-cases/go-1](https://www.alibabacloud.com/help/en/oss/use-cases/go-1)
