package apisign

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/url"
	"racent.com/pkg/config"
	"racent.com/pkg/helpers"
	"racent.com/pkg/logger"
	"racent.com/pkg/redis"
	"sort"
	"strconv"
	"time"
)

var (
	SignatureVersionV1       = "1.0"
	SignatureMethodMD5       = "md5"
	RedisPrefix              = "apisign:"
	RedisTll           int64 = 300
)

type Client struct {
	Gin          *gin.Context
	AccessKey    string
	AccessSecret string
}

func (c *Client) getCommonParams() map[string]string {
	params := map[string]string{
		"access_key":        c.AccessKey,
		"signature_nonce":   uuid.NewString(),
		"timestamp":         fmt.Sprintf("%d", time.Now().Unix()),
		"signature_version": SignatureVersionV1,
		"signature_method":  SignatureMethodMD5,
	}
	return params
}

func (c *Client) GenerateQueryParams(httpMethod string, queryParams map[string]string, body string) string {
	//queryMap, err := url.ParseQuery(queryParams)
	allKey := make([]string, 0)
	allVal := make(map[string]string)
	if !helpers.Empty(queryParams) {
		for i, i2 := range queryParams {
			allKey = append(allKey, i)
			allVal[i] = i2
		}
	}
	commonMap := c.getCommonParams()
	for j, j2 := range commonMap {
		allKey = append(allKey, j)
		allVal[j] = j2
	}

	logger.DebugJSON("GenerateQueryParams", "old_params", allKey)
	sort.Strings(allKey)
	logger.DebugJSON("GenerateQueryParams", "sort_params", allKey)

	stringToSign := generateStringToSign(allKey, allVal)
	logger.DebugString("GenerateQueryParams", "stringToSign", stringToSign)

	signature := c.generateSignature(httpMethod, stringToSign, body)

	return stringToSign + "&signature=" + signature
}

func (c *Client) generateSignature(httpMethod, stringToSign, body string) string {
	var signature string
	if helpers.Empty(body) {
		signature = helpers.MD5(c.AccessSecret + helpers.MD5(httpMethod+stringToSign))
	} else {
		signature = helpers.MD5(c.AccessSecret + helpers.MD5(httpMethod+stringToSign) + helpers.MD5(body))
	}
	return signature
}

func generateStringToSign(Keys []string, values map[string]string) string {
	var stringToSign string
	for k, v := range Keys {
		if k == 0 {
			stringToSign += url.PathEscape(v) + "=" + url.PathEscape(values[v])
		} else {
			stringToSign += "&" + url.PathEscape(v) + "=" + url.PathEscape(values[v])
		}
	}
	return stringToSign
}

func (c *Client) VerifySignature(httpMethod string, queryParams map[string]string, body string) error {
	// 验证公共参数
	err := c.verifyCommonParams(queryParams)
	if err != nil {
		return err
	}
	oldSign := queryParams["signature"]
	allKey := make([]string, 0)
	allVal := make(map[string]string)

	for k, v := range queryParams {
		if k == "signature" {
			continue
		}
		allKey = append(allKey, k)
		allVal[k] = v
	}
	//logger.DebugJSON("VerifySignature", "old_params", allKey)
	sort.Strings(allKey)
	//logger.DebugJSON("VerifySignature", "sort_params", allKey)

	stringToSign := generateStringToSign(allKey, allVal)
	//logger.DebugString("VerifySignature", "stringToSign", stringToSign)

	newSign := c.generateSignature(httpMethod, stringToSign, body)
	//logger.DebugString("VerifySignature", "new", newSign)
	if oldSign != newSign {
		return errors.New("signature 参数不合法")
	}

	return nil
}

func (c *Client) verifyCommonParams(queryParams map[string]string) error {

	err := c.verifyTimestamp(queryParams["timestamp"])
	if err != nil {
		return err
	}
	err = c.verifySignatureVersion(queryParams["signature_version"])
	if err != nil {
		return err
	}
	err = c.verifySignatureMethod(queryParams["signature_method"])
	if err != nil {
		return err
	}

	err = c.verifySignatureNonce(queryParams["signature_nonce"])
	if err != nil {
		return err
	}

	err = c.verifyAccessKey(queryParams["access_key"])
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) verifyAccessKey(accessKey string) error {
	if helpers.Empty(accessKey) {
		return errors.New("access_key 参数不能为空")
	}
	appName := config.GetString("app.name")
	if accessKey != appName {
		return errors.New("access_key 不合法")
	}
	c.AccessSecret = config.GetString("app.key")

	return nil
}

func (c *Client) verifySignatureNonce(signatureNonce string) error {
	if helpers.Empty(signatureNonce) {
		return errors.New("signature_nonce 参数不能为空")
	}
	lockKey := RedisPrefix + "nonce:" + signatureNonce
	exist, err := redis.Client.Get(lockKey)
	if err != nil {
		return err
	}
	if !helpers.Empty(exist) {
		return errors.New("signature_nonce 参数不合法")
	}
	_, err = redis.Client.Set(lockKey, 1, time.Duration(RedisTll)*time.Second)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) verifyTimestamp(t string) error {
	if helpers.Empty(t) {
		return errors.New("timestamp 参数不能为空")
	}

	t2, err := strconv.ParseInt(t, 10, 64)
	if err != nil {
		return err
	}
	nowTime := time.Now().Unix()
	if t2 < (nowTime-RedisTll) || t2 > (nowTime+RedisTll) {
		return errors.New("timestamp 参数不合法")
	}

	return nil
}

func (c *Client) verifySignatureVersion(version string) error {
	if helpers.Empty(version) {
		return errors.New("signature_version 参数不能为空")
	}
	if version != SignatureVersionV1 {
		return errors.New("signature_version 参数不合法")
	}
	return nil
}

func (c *Client) verifySignatureMethod(method string) error {
	if helpers.Empty(method) {
		return errors.New("signature_method 参数不能为空")
	}
	if method != SignatureMethodMD5 {
		return errors.New("signature_method 参数不合法")
	}
	return nil
}
