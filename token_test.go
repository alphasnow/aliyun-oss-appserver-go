package appserver

import (
	"encoding/json"
	"fmt"
	"sync"
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

	expectTokenStr := `{"OSSAccessKeyId":"yourAccessKeyId","policy":"eyJleHBpcmF0aW9uIjoiMjAyNS0wMS0wMVQwMDowMDowMFoiLCJjb25kaXRpb25zIjpbWyJzdGFydHMtd2l0aCIsIiRrZXkiLCJ1c2VyLWRpci1wcmVmaXgvIl0seyJidWNrZXQiOiJidWNrZXQtbmFtZSJ9LFsiY29udGVudC1sZW5ndGgtcmFuZ2UiLDEsMTA0ODU3NjBdLFsiaW4iLCIkY29udGVudC10eXBlIixbImltYWdlL2pwZWciLCJpbWFnZS9wbmciXV1dfQ==","callback":"eyJjYWxsYmFja1VybCI6Imh0dHA6Ly9kb21haW4uY29tL29zcy9jYWxsYmFjayIsImNhbGxiYWNrQm9keSI6IntcImJ1Y2tldFwiOiR7YnVja2V0fSxcIm9iamVjdFwiOiR7b2JqZWN0fSxcImV0YWdcIjoke2V0YWd9LFwic2l6ZVwiOiR7c2l6ZX0sXCJtaW1lVHlwZVwiOiR7bWltZVR5cGV9LFwiaW1hZ2VJbmZvXCI6e1wiaGVpZ2h0XCI6JHtpbWFnZUluZm8uaGVpZ2h0fSxcIndpZHRoXCI6JHtpbWFnZUluZm8ud2lkdGh9LFwiZm9ybWF0XCI6JHtpbWFnZUluZm8uZm9ybWF0fX0sXCJjcmM2NFwiOiR7Y3JjNjR9LFwiY29udGVudE1kNVwiOiR7Y29udGVudE1kNX0sXCJ2cGNJZFwiOiR7dnBjSWR9LFwiY2xpZW50SXBcIjoke2NsaWVudElwfSxcInJlcUlkXCI6JHtyZXFJZH0sXCJvcGVyYXRpb25cIjoke29wZXJhdGlvbn19IiwiY2FsbGJhY2tCb2R5VHlwZSI6ImFwcGxpY2F0aW9uL2pzb24ifQ==","signature":"wQZPtbuNzqTOol/oXZHIv7SLhc0=","host":"https://bucket-name.oss-cn-hangzhou.aliyuncs.com","expire":1735689600,"directory":"user-dir-prefix/"}`
	if tokenJsonStr != expectTokenStr {
		t.Error("token error")
	}
}

func TestTokenPolicy(t *testing.T) {

	token := NewToken(&Config{
		AccessKeyId:     "yourAccessKeyId",
		AccessKeySecret: "yourAccessKeySecret",
		Host:            "https://bucket-name.oss-cn-hangzhou.aliyuncs.com",
		Directory:       "user-dir-prefix/",
		ExpireSecond:    600,
		CallbackUrl:     "http://domain.com/oss/callback",
	})

	g := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		g.Add(1)
		go func(idx int) {
			defer g.Done()
			dir := fmt.Sprintf("user-%d-prefix/", idx)
			policy := new(Policy)
			policy.SetDirectory(dir)
			tokenPayload, _ := token.SetPolicy(policy).Generate()
			if tokenPayload.Directory != dir {
				t.Error("dir error", dir)
			}
		}(i)
	}
	g.Wait()
}
