package appserver

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"hash"
	"io"
	"time"
)

const TimeGMTISO8601 = "2006-01-02T15:04:05Z"
const CallbackBodyParam = `{"bucket":${bucket},"object":${object},"etag":${etag},"size":${size},"mimeType":${mimeType},"imageInfo":{"height":${imageInfo.height},"width":${imageInfo.width},"format":${imageInfo.format}},"crc64":${crc64},"contentMd5":${contentMd5},"vpcId":${vpcId},"clientIp":${clientIp},"reqId":${reqId},"operation":${operation}}`
const CallbackBodyTypeParam = "application/json"
const DefaultExpireSecond = 600

type Token struct {
	config   *Config
	policy   *Policy
	callback *Callback
}

func NewToken(config *Config) *Token {
	var callback *Callback
	if config.CallbackUrl != "" {
		callback = newCallback(config)
	}
	return &Token{
		config:   config,
		callback: callback,
	}
}

func newPolicy(config *Config) *Policy {
	sp := new(Policy)
	var expireSecond int64
	if config.ExpireSecond != 0 {
		expireSecond = config.ExpireSecond
	} else {
		expireSecond = DefaultExpireSecond
	}
	expireTime := time.Now().Add(time.Duration(expireSecond) * time.Second)
	sp.SetExpireTime(expireTime)
	if config.Directory != "" {
		sp.SetDirectory(config.Directory)
	}
	return sp
}

func newCallback(config *Config) *Callback {
	cp := new(Callback)
	cp.CallbackUrl = config.CallbackUrl

	if config.CallbackBody != "" {
		cp.CallbackBody = config.CallbackBody
	} else {
		cp.CallbackBody = CallbackBodyParam
	}

	if config.CallbackBodyType != "" {
		cp.CallbackBodyType = config.CallbackBodyType
	} else {
		cp.CallbackBodyType = CallbackBodyTypeParam
	}
	return cp
}

func (t *Token) SetPolicy(policy *Policy) *Token {
	k := *t
	k.policy = policy
	return &k
}

func (t *Token) SetCallback(callback *Callback) *Token {
	k := *t
	k.callback = callback
	return &k
}

func (t *Token) Generate() (*SignatureToken, error) {
	// policy
	if t.policy == nil {
		t.policy = newPolicy(t.config)
	}
	policyByte, err := json.Marshal(t.policy)
	policy := base64.StdEncoding.EncodeToString(policyByte)

	// signature
	h := hmac.New(func() hash.Hash { return sha1.New() }, []byte(t.config.AccessKeySecret))
	_, err = io.WriteString(h, policy)
	if err != nil {
		return nil, err
	}
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))

	// callback
	var callback string
	if t.callback.isValid() {
		var callbackStr []byte
		callbackStr, err = json.Marshal(t.callback)
		if err != nil {
			return nil, err
		}
		callback = base64.StdEncoding.EncodeToString(callbackStr)
	}

	// token
	var policyToken SignatureToken
	policyToken.OSSAccessKeyId = t.config.AccessKeyId
	policyToken.Host = t.config.Host
	policyToken.Directory = t.policy.GetDirectory()
	policyToken.Expire = t.policy.GetExpire()
	policyToken.Signature = signature
	policyToken.Policy = policy
	policyToken.Callback = callback

	return &policyToken, nil
}

type Config struct {
	// SignatureToken
	AccessKeyId     string `json:"access_key_id"`
	AccessKeySecret string `json:"access_key_secret"`
	Host            string `json:"host"`

	// Callback
	CallbackUrl      string `json:"callback_url"`
	CallbackBody     string `json:"callback_body"`
	CallbackBodyType string `json:"callback_body_type"`

	// Policy
	Directory    string `json:"directory"`
	ExpireSecond int64  `json:"expire_second"`
}

// Policy
// https://help.aliyun.com/zh/oss/developer-reference/signature-version-1
type Policy struct {
	expiredAt time.Time
	uploadDir string

	Expiration string `json:"expiration"` // required
	Conditions []any  `json:"conditions"` // optional
}

func (c *Policy) GetExpire() int64 {
	return c.expiredAt.Unix()
}

func (c *Policy) GetDirectory() string {
	return c.uploadDir
}

func (c *Policy) SetExpireTime(expiredAt time.Time) {
	c.Expiration = expiredAt.UTC().Format(TimeGMTISO8601)
	c.expiredAt = expiredAt
}

func (c *Policy) SetDirectory(uploadDir string) {
	c.Conditions = append(c.Conditions, []string{
		"starts-with", "$key", uploadDir,
	})

	c.uploadDir = uploadDir
}

func (c *Policy) SetBucket(bucket string) {
	c.Conditions = append(c.Conditions, map[string]string{
		"bucket": bucket,
	})
}

func (c *Policy) SetContentLengthRange(min uint64, max uint64) {
	c.Conditions = append(c.Conditions, []any{
		"content-length-range", min, max,
	})
}

func (c *Policy) SetContentType(types ...string) {
	c.Conditions = append(c.Conditions, []any{
		"in", "$content-type", types,
	})
}

// SignatureToken
// https://help.aliyun.com/zh/oss/developer-reference/postobject
type SignatureToken struct {
	// post object param
	OSSAccessKeyId string `json:"OSSAccessKeyId"` // required
	Policy         string `json:"policy"`         // required
	Callback       string `json:"callback"`       // optional
	Signature      string `json:"signature"`      // required
	// api param
	Host      string `json:"host"`      // optional
	Expire    int64  `json:"expire"`    // optional
	Directory string `json:"directory"` // optional

}

// Callback
// https://help.aliyun.com/zh/oss/developer-reference/callback
type Callback struct {
	CallbackUrl      string `json:"callbackUrl"`                // required
	CallbackBody     string `json:"callbackBody"`               // optional
	CallbackBodyType string `json:"callbackBodyType,omitempty"` // optional, default: application/x-www-form-urlencoded
}

func (c *Callback) isValid() bool {
	return c.CallbackUrl != "" && c.CallbackBody != ""
}
