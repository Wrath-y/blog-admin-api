package admin

import (
	"blog-admin-api/core"
	"blog-admin-api/errcode"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"hash"
	"time"
)

type ConfigStruct struct {
	Expiration string     `json:"expiration"`
	Conditions [][]string `json:"conditions"`
}

type PolicyToken struct {
	AccessKeyId string `json:"accessid"`
	Host        string `json:"host"`
	Expire      int64  `json:"expire"`
	Signature   string `json:"signature"`
	Policy      string `json:"policy"`
	Directory   string `json:"dir"`
	Callback    string `json:"callback"`
}

type Data struct {
	AccessUrl string      `json:"access_url"`
	Drive     string      `json:"drive"`
	FileField string      `json:"file_field"`
	Form      interface{} `json:"form"`
	Headers   interface{} `json:"headers"`
	UploadUrl string      `json:"upload_url"`
}

type Form struct {
	OSSAccessKeyId      string `json:"OSSAccessKeyId"`
	Policy              string `json:"policy"`
	Signature           string `json:"signature"`
	SuccessActionStatus int    `json:"success_action_status"`
}

type CallbackParam struct {
	CallbackUrl      string `json:"callbackUrl"`
	CallbackBody     string `json:"callbackBody"`
	CallbackBodyType string `json:"callbackBodyType"`
}

func GetGmtIso8601(expireEnd int64) string {
	var tokenExpire = time.Unix(expireEnd, 0).Format("2006-01-02T15:04:05Z")
	return tokenExpire
}

func GetUpload(c *core.Context) {
	logMap := make(map[string]interface{})
	accessKeyId := viper.GetString("aliyun.oss.access_keyid")
	accessKeySecret := viper.GetString("aliyun.oss.access_keysecret")
	// host的格式为 bucketname.endpoint
	host := fmt.Sprintf("https://%s.oss-cn-shanghai.aliyuncs.com", viper.GetString("aliyun.oss.bucket"))
	logMap["host"] = host
	logMap["accessKeyId"] = accessKeyId
	logMap["accessKeySecret"] = accessKeySecret

	// 上传文件时指定的前缀。
	uploadDir := ""
	expireTime := 30
	now := time.Now().Unix()
	expireEnd := now + int64(expireTime)
	var tokenExpire = GetGmtIso8601(expireEnd)

	//create post policy json
	var config ConfigStruct
	config.Expiration = tokenExpire
	var condition []string
	condition = append(condition, "starts-with")
	condition = append(condition, "$key")
	condition = append(condition, uploadDir)
	config.Conditions = append(config.Conditions, condition)
	logMap["config"] = config

	//calucate signature
	result, err := json.Marshal(config)
	if err != nil {
		c.ErrorL("序列化config失败", logMap, err.Error())
		c.FailWithErrCode(errcode.AdminNetworkBusy, nil)
		return
	}
	debyte := base64.StdEncoding.EncodeToString(result)
	h := hmac.New(func() hash.Hash { return sha1.New() }, []byte(accessKeySecret))
	h.Write([]byte(debyte))
	signedStr := base64.StdEncoding.EncodeToString(h.Sum(nil))

	var policyToken PolicyToken
	policyToken.AccessKeyId = accessKeyId
	policyToken.Host = host
	policyToken.Expire = expireEnd
	policyToken.Signature = signedStr
	policyToken.Directory = uploadDir
	policyToken.Policy = debyte
	logMap["policyToken"] = policyToken

	response := &Data{
		AccessUrl: policyToken.Host,
		Drive:     "oss",
		FileField: "file",
		Form: &Form{
			OSSAccessKeyId:      accessKeyId,
			Policy:              policyToken.Policy,
			Signature:           policyToken.Signature,
			SuccessActionStatus: 200,
		},
		Headers:   []int{},
		UploadUrl: policyToken.Host,
	}

	c.Success(response)
}
