package appserver

import (
	"encoding/json"
	"testing"
	"time"
)

func TestTokenGenerate(t *testing.T) {

	token := NewToken(&Config{
		AccessKeyId:     "yourAccessKeyId",
		AccessKeySecret: "yourAccessKeySecret",
		Host:            "https://bucket-name.oss-cn-hangzhou.aliyuncs.com",
		ExpireSecond:    600,
	})

	targetTime, _ := time.Parse("2006-01-02 15:04:05", "2025-01-01 00:00:00")
	policy := new(Policy)
	policy.SetExpireTime(targetTime)

	tokenPayload, _ := token.SetPolicy(policy).Generate()
	tokenJson, _ := json.Marshal(tokenPayload)
	tokenJsonStr := string(tokenJson)

	expectTokenStr := `{"OSSAccessKeyId":"yourAccessKeyId","policy":"eyJleHBpcmF0aW9uIjoiMjAyNS0wMS0wMVQwMDowMDowMFoiLCJjb25kaXRpb25zIjpudWxsfQ==","callback":"","signature":"S7QSuk+DEd0QdMRZFhwv3yjuE6g=","host":"https://bucket-name.oss-cn-hangzhou.aliyuncs.com","expire":1735689600,"directory":""}`
	if tokenJsonStr != expectTokenStr {
		t.Error("token error")
	}
}

func TestTokenWithConfigGenerate(t *testing.T) {

	token := NewToken(&Config{
		AccessKeyId:     "yourAccessKeyId",
		AccessKeySecret: "yourAccessKeySecret",
		Host:            "https://bucket-name.oss-cn-hangzhou.aliyuncs.com",
		Directory:       "user-dir-prefix/",
		ExpireSecond:    600,
		CallbackUrl:     "http://domain.com/oss/callback",
	})

	targetTime, _ := time.Parse("2006-01-02 15:04:05", "2025-01-01 00:00:00")
	policy := new(Policy)
	policy.SetExpireTime(targetTime)
	policy.SetDirectory("user-dir-prefix/")

	tokenPayload, _ := token.SetPolicy(policy).Generate()
	tokenJson, _ := json.Marshal(tokenPayload)
	tokenJsonStr := string(tokenJson)

	expectTokenStr := `{"OSSAccessKeyId":"yourAccessKeyId","policy":"eyJleHBpcmF0aW9uIjoiMjAyNS0wMS0wMVQwMDowMDowMFoiLCJjb25kaXRpb25zIjpbWyJzdGFydHMtd2l0aCIsIiRrZXkiLCJ1c2VyLWRpci1wcmVmaXgvIl1dfQ==","callback":"eyJjYWxsYmFja1VybCI6Imh0dHA6Ly9kb21haW4uY29tL29zcy9jYWxsYmFjayIsImNhbGxiYWNrQm9keSI6IntcImJ1Y2tldFwiOiR7YnVja2V0fSxcIm9iamVjdFwiOiR7b2JqZWN0fSxcImV0YWdcIjoke2V0YWd9LFwic2l6ZVwiOiR7c2l6ZX0sXCJtaW1lVHlwZVwiOiR7bWltZVR5cGV9LFwiaW1hZ2VJbmZvXCI6e1wiaGVpZ2h0XCI6JHtpbWFnZUluZm8uaGVpZ2h0fSxcIndpZHRoXCI6JHtpbWFnZUluZm8ud2lkdGh9LFwiZm9ybWF0XCI6JHtpbWFnZUluZm8uZm9ybWF0fX0sXCJjcmM2NFwiOiR7Y3JjNjR9LFwiY29udGVudE1kNVwiOiR7Y29udGVudE1kNX0sXCJ2cGNJZFwiOiR7dnBjSWR9LFwiY2xpZW50SXBcIjoke2NsaWVudElwfSxcInJlcUlkXCI6JHtyZXFJZH0sXCJvcGVyYXRpb25cIjoke29wZXJhdGlvbn19IiwiY2FsbGJhY2tCb2R5VHlwZSI6ImFwcGxpY2F0aW9uL2pzb24ifQ==","signature":"uXL82wU5IGCd7vcZKX9gua5TUJs=","host":"https://bucket-name.oss-cn-hangzhou.aliyuncs.com","expire":1735689600,"directory":"user-dir-prefix/"}`
	if tokenJsonStr != expectTokenStr {
		t.Error("token error")
	}
	//{
	//    "OSSAccessKeyId": "yourAccessKeyId",
	//    "policy": "eyJleHBpcmF0aW9uIjoiMjAyNS0wMS0wMVQwMDowMDowMFoiLCJjb25kaXRpb25zIjpbWyJzdGFydHMtd2l0aCIsIiRrZXkiLCJ1c2VyLWRpci1wcmVmaXgvIl1dfQ==",
	//    "callback": "eyJjYWxsYmFja1VybCI6Imh0dHA6Ly9kb21haW4uY29tL29zcy9jYWxsYmFjayIsImNhbGxiYWNrQm9keSI6IntcImJ1Y2tldFwiOiR7YnVja2V0fSxcIm9iamVjdFwiOiR7b2JqZWN0fSxcImV0YWdcIjoke2V0YWd9LFwic2l6ZVwiOiR7c2l6ZX0sXCJtaW1lVHlwZVwiOiR7bWltZVR5cGV9LFwiaW1hZ2VJbmZvXCI6e1wiaGVpZ2h0XCI6JHtpbWFnZUluZm8uaGVpZ2h0fSxcIndpZHRoXCI6JHtpbWFnZUluZm8ud2lkdGh9LFwiZm9ybWF0XCI6JHtpbWFnZUluZm8uZm9ybWF0fX0sXCJjcmM2NFwiOiR7Y3JjNjR9LFwiY29udGVudE1kNVwiOiR7Y29udGVudE1kNX0sXCJ2cGNJZFwiOiR7dnBjSWR9LFwiY2xpZW50SXBcIjoke2NsaWVudElwfSxcInJlcUlkXCI6JHtyZXFJZH0sXCJvcGVyYXRpb25cIjoke29wZXJhdGlvbn19IiwiY2FsbGJhY2tCb2R5VHlwZSI6ImFwcGxpY2F0aW9uL2pzb24ifQ==",
	//    "signature": "uXL82wU5IGCd7vcZKX9gua5TUJs=",
	//    "host": "https://bucket-name.oss-cn-hangzhou.aliyuncs.com",
	//    "expire": 1735689600,
	//    "directory": "user-dir-prefix/"
	//}
}

func TestTokenWithPolicyGenerate(t *testing.T) {

	token := NewToken(&Config{
		AccessKeyId:     "yourAccessKeyId",
		AccessKeySecret: "yourAccessKeySecret",
		Host:            "https://bucket-name.oss-cn-hangzhou.aliyuncs.com",
		Directory:       "user-dir-prefix/",
		ExpireSecond:    600,
		CallbackUrl:     "http://domain.com/oss/callback",
	})

	targetTime, _ := time.Parse("2006-01-02 15:04:05", "2025-01-01 00:00:00")
	policy := new(Policy)
	policy.SetExpireTime(targetTime)
	policy.SetDirectory("user-dir-prefix/")
	policy.SetBucket("bucket-name")
	policy.SetContentLengthRange(1, 10*1024*1024)
	policy.SetContentType("image/jpeg", "image/png")

	tokenPayload, _ := token.SetPolicy(policy).Generate()
	tokenJson, _ := json.Marshal(tokenPayload)
	tokenJsonStr := string(tokenJson)

	expectTokenStr := `{"OSSAccessKeyId":"yourAccessKeyId","policy":"eyJleHBpcmF0aW9uIjoiMjAyNS0wMS0wMVQwMDowMDowMFoiLCJjb25kaXRpb25zIjpbWyJzdGFydHMtd2l0aCIsIiRrZXkiLCJ1c2VyLWRpci1wcmVmaXgvIl0seyJidWNrZXQiOiJidWNrZXQtbmFtZSJ9LFsiY29udGVudC1sZW5ndGgtcmFuZ2UiLDEsMTA0ODU3NjBdLFsiaW4iLCIkY29udGVudC10eXBlIixbImltYWdlL2pwZWciLCJpbWFnZS9wbmciXV1dfQ==","callback":"eyJjYWxsYmFja1VybCI6Imh0dHA6Ly84OC44OC44OC44ODo4ODg4IiwiY2FsbGJhY2tCb2R5Ijoie1wiYnVja2V0XCI6JHtidWNrZXR9LFwib2JqZWN0XCI6JHtvYmplY3R9LFwiZXRhZ1wiOiR7ZXRhZ30sXCJzaXplXCI6JHtzaXplfSxcIm1pbWVUeXBlXCI6JHttaW1lVHlwZX0sXCJpbWFnZUluZm9cIjp7XCJoZWlnaHRcIjoke2ltYWdlSW5mby5oZWlnaHR9LFwid2lkdGhcIjoke2ltYWdlSW5mby53aWR0aH0sXCJmb3JtYXRcIjoke2ltYWdlSW5mby5mb3JtYXR9fSxcImNyYzY0XCI6JHtjcmM2NH0sXCJjb250ZW50TWQ1XCI6JHtjb250ZW50TWQ1fSxcInZwY0lkXCI6JHt2cGNJZH0sXCJjbGllbnRJcFwiOiR7Y2xpZW50SXB9LFwicmVxSWRcIjoke3JlcUlkfSxcIm9wZXJhdGlvblwiOiR7b3BlcmF0aW9ufX0iLCJjYWxsYmFja0JvZHlUeXBlIjoiYXBwbGljYXRpb24vanNvbiJ9","signature":"wQZPtbuNzqTOol/oXZHIv7SLhc0=","host":"https://bucket-name.oss-cn-hangzhou.aliyuncs.com","expire":1735689600,"directory":"user-dir-prefix/"}`
	if tokenJsonStr != expectTokenStr {
		t.Error("token error")
	}
}
