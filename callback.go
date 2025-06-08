package appserver

import (
	"crypto"
	"crypto/md5"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// CallbackBody 结构体用于存储文件相关信息以及请求相关的一些元数据
type CallbackBody struct {
	// Bucket是存储空间名称
	Bucket string `json:"bucket"`
	// Object表示对象（文件）的完整路径
	Object string `json:"object"`
	// Etag是文件的ETag，即返回给用户的ETag字段
	Etag string `json:"etag"`
	// Size代表Object大小，调用CompleteMultipartUpload时，size为整个Object的大小
	Size int `json:"size"`
	// MimeType代表资源类型，例如jpeg图片的资源类型为image/jpeg
	MimeType string `json:"mimeType"`
	// ImageInfo用于存储图片相关的额外信息，仅适用于图片格式
	ImageInfo ImageInfo `json:"imageInfo"`
	// Crc64与上传文件后返回的x-oss-hash-crc64ecma头内容一致
	Crc64 uint `json:"crc64"`
	// ContentMd5与上传文件后返回的Content-MD5头内容一致，仅在调用PutObject和PostObject接口上传文件时，该变量的值不为空
	ContentMd5 string `json:"contentMd5"`
	// VpcId发起请求的客户端所在的VpcId，如果不是通过VPC发起请求，则该变量的值为空
	VpcId *string `json:"vpcId"`
	// ClientIp发起请求的客户端IP地址
	ClientIp string `json:"clientIp"`
	// ReqId发起请求的RequestId
	ReqId string `json:"reqId"`
	// Operation发起请求的接口名称，例如PutObject、PostObject等
	Operation string `json:"operation"`
}

type ImageInfo struct {
	// Height是图片高度，对于非图片格式，该变量的值为空
	Height int `json:"height"`
	// Width是图片宽度，对于非图片格式，该变量的值为空
	Width int `json:"width"`
	// Format是图片格式，例如JPG、PNG等，对于非图片格式，该变量的值为空
	Format string `json:"format"`
}

const PubKeyUrlHeader = "X-Oss-Pub-Key-Url"
const AuthorizationHeader = "Authorization"

type AliyunOSSCallback struct {
	req *http.Request
}

func NewAliyunOSSCallback(req *http.Request) *AliyunOSSCallback {
	return &AliyunOSSCallback{req: req}
}

func (a *AliyunOSSCallback) VerifySignature() (*CallbackBody, error) {
	bodyContent, err := io.ReadAll(a.req.Body)
	if err != nil {
		return nil, err
	}
	defer a.req.Body.Close()
	byteMd5, err := GetMD5FromNewAuthString(bodyContent, a.req.URL.Path, a.req.URL.RawQuery)
	if err != nil {
		return nil, err
	}

	publicKeyURLBase64 := a.req.Header.Get(PubKeyUrlHeader)
	bytePublicKey, err := GetPublicKey(publicKeyURLBase64)
	if err != nil {
		return nil, err
	}

	strAuthorizationBase64 := a.req.Header.Get(AuthorizationHeader)
	authorization, err := GetAuthorization(strAuthorizationBase64)
	if err != nil {
		return nil, err
	}

	if err = VerifySignature(bytePublicKey, byteMd5, authorization); err != nil {
		return nil, err
	}

	callbackBody := new(CallbackBody)
	if err = json.Unmarshal(bodyContent, callbackBody); err != nil {
		return nil, err
	}
	return callbackBody, nil
}

func VerifySignature(bytePublicKey []byte, byteMd5 []byte, authorization []byte) error {
	pubBlock, _ := pem.Decode(bytePublicKey)
	if pubBlock == nil {
		// fmt.Printf("Failed to parse PEM block containing the public key")
		return errors.New("failed to parse PEM block containing the public key")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(pubBlock.Bytes)
	if (pubInterface == nil) || (err != nil) {
		//fmt.Printf("x509.ParsePKIXPublicKey(publicKey) failed : %s \n", err.Error())
		return fmt.Errorf("x509.ParsePKIXPublicKey(publicKey) failed : %w \n", err)
	}
	pub := pubInterface.(*rsa.PublicKey)

	errorVerifyPKCS1v15 := rsa.VerifyPKCS1v15(pub, crypto.MD5, byteMd5, authorization)
	if errorVerifyPKCS1v15 != nil {
		//fmt.Printf("\nSignature Verification is Failed : %s \n", errorVerifyPKCS1v15.Error())
		//printByteArray(byteMd5, "AuthMd5(fromNewAuthString)")
		//printByteArray(bytePublicKey, "PublicKeyBase64")
		//printByteArray(authorization, "AuthorizationFromRequest")
		return fmt.Errorf("Signature Verification is Failed : %v \n", errorVerifyPKCS1v15.Error())
	}

	// fmt.Printf("\nSignature Verification is Successful. \n")
	return nil
}

// GetPublicKey : Get PublicKey bytes from Request.URL
func GetPublicKey(publicKeyURLBase64 string) ([]byte, error) {
	var bytePublicKey []byte
	publicKeyURL, err := base64.StdEncoding.DecodeString(publicKeyURLBase64)
	if err != nil {
		return nil, err
	}
	// fmt.Printf("publicKeyURL={%s}\n", publicKeyURL)

	// get PublicKey Content from URL
	responsePublicKeyURL, err := http.Get(string(publicKeyURL))
	if err != nil {
		// fmt.Printf("Get PublicKey Content from URL failed : %s \n", err.Error())
		return nil, err
	}
	bytePublicKey, err = io.ReadAll(responsePublicKeyURL.Body)
	if err != nil {
		// fmt.Printf("Read PublicKey Content from URL failed : %s \n", err.Error())
		return bytePublicKey, err
	}
	defer responsePublicKeyURL.Body.Close()

	// fmt.Printf("publicKey={%s}\n", bytePublicKey)
	return bytePublicKey, nil
}

// GetAuthorization : decode from Base64String
func GetAuthorization(strAuthorizationBase64 string) ([]byte, error) {
	var byteAuthorization []byte
	var err error
	// Get Authorization bytes : decode from Base64String
	if strAuthorizationBase64 == "" {
		//fmt.Println("Failed to get authorization field from request header. ")
		return nil, errors.New("no authorization field in Request header")
	}

	byteAuthorization, err = base64.StdEncoding.DecodeString(strAuthorizationBase64)
	if err != nil {
		return nil, err
	}
	return byteAuthorization, nil
}

// GetMD5FromNewAuthString : Get MD5 bytes from Newly Constructed Authorization String.
func GetMD5FromNewAuthString(bodyContent []byte, urlPath string, urlQuery string) ([]byte, error) {
	var byteMD5 []byte
	// Construct the New Auth String from URI+Query+Body
	strCallbackBody := string(bodyContent)
	// fmt.Printf("r.URL.RawPath={%s}, r.URL.Query()={%s}, strCallbackBody={%s}\n", r.URL.RawPath, r.URL.Query(), strCallbackBody)
	strURLPathDecode, errUnescape := unescapePath(urlPath, encodePathSegment) //url.PathUnescape(r.URL.Path) for Golang v1.8.2+
	if errUnescape != nil {
		// fmt.Printf("url.PathUnescape failed : URL.Path=%s, error=%s \n", r.URL.Path, err.Error())
		return nil, errUnescape
	}

	// Generate New Auth String prepare for MD5
	strAuth := ""
	if urlQuery == "" {
		strAuth = fmt.Sprintf("%s\n%s", strURLPathDecode, strCallbackBody)
	} else {
		strAuth = fmt.Sprintf("%s?%s\n%s", strURLPathDecode, urlPath, strCallbackBody)
	}
	// fmt.Printf("NewlyConstructedAuthString={%s}\n", strAuth)

	// Generate MD5 from the New Auth String
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(strAuth))
	byteMD5 = md5Ctx.Sum(nil)

	return byteMD5, nil
}
