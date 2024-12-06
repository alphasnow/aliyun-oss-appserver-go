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
		Directory:       "user-dir-prefix/",
		ExpireSecond:    30,
		CallbackUrl:     "http://88.88.88.88:8888",
	})

	targetTime, _ := time.Parse("2006-01-02 15:04:05", "2025-01-01 00:00:00")
	policy := new(Policy)
	policy.SetDirectory("user-dir-prefix/")
	policy.SetExpireTime(targetTime)

	tokenPayload, _ := token.SetPolicy(policy).Generate()
	tokenJson, _ := json.Marshal(tokenPayload)
	tokenJsonStr := string(tokenJson)

	expectTokenStr := `{"accessid":"yourAccessKeyId","host":"https://bucket-name.oss-cn-hangzhou.aliyuncs.com","expire":1735689600,"signature":"uXL82wU5IGCd7vcZKX9gua5TUJs=","policy":"eyJleHBpcmF0aW9uIjoiMjAyNS0wMS0wMVQwMDowMDowMFoiLCJjb25kaXRpb25zIjpbWyJzdGFydHMtd2l0aCIsIiRrZXkiLCJ1c2VyLWRpci1wcmVmaXgvIl1dfQ==","dir":"user-dir-prefix/","callback":"eyJjYWxsYmFja1VybCI6Imh0dHA6Ly84OC44OC44OC44ODo4ODg4IiwiY2FsbGJhY2tCb2R5Ijoie1wiYnVja2V0XCI6JHtidWNrZXR9LFwib2JqZWN0XCI6JHtvYmplY3R9LFwiZXRhZ1wiOiR7ZXRhZ30sXCJzaXplXCI6JHtzaXplfSxcIm1pbWVUeXBlXCI6JHttaW1lVHlwZX0sXCJpbWFnZUluZm9cIjp7XCJoZWlnaHRcIjoke2ltYWdlSW5mby5oZWlnaHR9LFwid2lkdGhcIjoke2ltYWdlSW5mby53aWR0aH0sXCJmb3JtYXRcIjoke2ltYWdlSW5mby5mb3JtYXR9fSxcImNyYzY0XCI6JHtjcmM2NH0sXCJjb250ZW50TWQ1XCI6JHtjb250ZW50TWQ1fSxcInZwY0lkXCI6JHt2cGNJZH0sXCJjbGllbnRJcFwiOiR7Y2xpZW50SXB9LFwicmVxSWRcIjoke3JlcUlkfSxcIm9wZXJhdGlvblwiOiR7b3BlcmF0aW9ufX0iLCJjYWxsYmFja0JvZHlUeXBlIjoiYXBwbGljYXRpb24vanNvbiJ9"}`
	if tokenJsonStr != expectTokenStr {
		t.Error("token error")
	}

}
