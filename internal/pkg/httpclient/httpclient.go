package httpclient

// http-client

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"racent.com/pkg/helpers"
	"racent.com/pkg/logger"
	"strings"
	"time"
)

type ResponseWrapper struct {
	StatusCode int
	Body       []byte
	Header     http.Header
}

type RequestData struct {
	Headers map[string]string
	Url     string
	Params  interface{}
	//Params     map[string]interface{}
	Cert       string
	PrivateKey string
	Files      []UploadFile
}

type UploadFile struct {
	Name     string
	Filepath string
}

var httpClient = &http.Client{Timeout: time.Second * 30}

func (request *RequestData) Get() ResponseWrapper {
	urlParams := url.Values{}
	Url, err := url.Parse(request.Url)
	if err != nil {
		return createRequestError(err)
	}
	if !helpers.Empty(request.Params) {
		for k, v := range request.Params.(map[string]interface{}) {
			urlParams.Set(k, fmt.Sprint(v))
		}
		Url.RawQuery = urlParams.Encode()
	}
	urlPath := Url.String()
	httpRequest, err := http.NewRequest("GET", urlPath, nil)
	if err != nil {
		return createRequestError(err)
	}
	err = addCerts(request)
	if err != nil {
		return createRequestError(err)
	}
	setRequestHeaders(httpRequest, request.Headers)
	//初始化httpClient
	return toRequest(httpRequest)
}

func (request *RequestData) Post() ResponseWrapper {
	logger.DebugString("httpclient", "request_url", request.Url)
	requestBody, realContentType, err := getReader(request)
	if err != nil {
		return createRequestError(err)
	}
	httpRequest, err := http.NewRequest("POST", request.Url, requestBody)
	if err != nil {
		return createRequestError(err)
	}
	err = addCerts(request)
	if err != nil {
		return createRequestError(err)
	}
	headers := request.Headers
	if headers == nil {
		headers = make(map[string]string)
	}
	headers["Content-Type"] = realContentType
	setRequestHeaders(httpRequest, headers)
	return toRequest(httpRequest)
}

func (request *RequestData) Put() ResponseWrapper {
	requestBody, realContentType, err := getReader(request)
	if err != nil {
		return createRequestError(err)
	}
	httpRequest, err := http.NewRequest("PUT", request.Url, requestBody)
	if err != nil {
		return createRequestError(err)
	}
	err = addCerts(request)
	if err != nil {
		return createRequestError(err)
	}

	headers := request.Headers
	if headers == nil {
		headers = make(map[string]string)
	}

	headers["Content-Type"] = realContentType
	setRequestHeaders(httpRequest, headers)
	return toRequest(httpRequest)
}

func (request *RequestData) Delete() ResponseWrapper {
	httpRequest, err := http.NewRequest("DELETE", request.Url, nil)
	if err != nil {
		return createRequestError(err)
	}
	err = addCerts(request)
	if err != nil {
		return createRequestError(err)
	}
	setRequestHeaders(httpRequest, request.Headers)
	return toRequest(httpRequest)
}

func setRequestHeaders(httpRequest *http.Request, headers map[string]string) {
	if headers != nil {
		for k, v := range headers {
			httpRequest.Header.Set(k, v)
		}
	}
}

func getReader(request *RequestData) (io.Reader, string, error) {
	if request.Files != nil {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		for _, uploadFile := range request.Files {
			file, err := os.Open(uploadFile.Filepath)
			if err != nil {
				return nil, "", err
			}
			part, err := writer.CreateFormFile(uploadFile.Name, filepath.Base(uploadFile.Filepath))
			if err != nil {
				return nil, "", err
			}
			_, err = io.Copy(part, file)
			if err != nil {
				return nil, "", err
			}
			file.Close()
		}
		for k, v := range request.Params.(map[string]interface{}) {
			if err := writer.WriteField(k, fmt.Sprint(v)); err != nil {
				return nil, "", err
			}
		}
		if err := writer.Close(); err != nil {
			return nil, "", err
		}
		return body, writer.FormDataContentType(), nil
	} else {
		headers := request.Headers
		contentType := headers["Content-Type"]
		switch contentType {
		case "application/x-www-form-urlencoded":
			urlValues := url.Values{}
			for key, val := range request.Params.(map[string]interface{}) {
				urlValues.Set(key, fmt.Sprint(val))
			}
			reqBody := urlValues.Encode()
			return strings.NewReader(reqBody), contentType, nil
		case "application/json":
			if helpers.Empty(request.Params) {
				return nil, "application/json", nil
			}
			logger.DebugJSON("httpclient", "request_header", request.Headers)
			logger.DebugString("httpclient", "request_body", string(request.Params.([]byte)))
			return bytes.NewReader(request.Params.([]byte)), "application/json", nil
		}
		return nil, "", errors.New("Content-Type 不合法")
	}
}

func toRequest(httpRequest *http.Request) ResponseWrapper {
	wrapper := ResponseWrapper{StatusCode: 0, Body: nil, Header: make(http.Header)}
	resp, err := httpClient.Do(httpRequest)

	if err != nil {
		wrapper.Body = []byte(fmt.Sprintf("执行HTTP请求错误-%s", err.Error()))
		logger.WarnJSON("httpclient", "resp", wrapper)
		return wrapper
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		wrapper.Body = []byte(fmt.Sprintf("读取HTTP请求返回值失败-%s", err.Error()))
		logger.WarnJSON("httpclient", "resp", wrapper)
		return wrapper
	}
	wrapper.StatusCode = resp.StatusCode
	wrapper.Body = body
	wrapper.Header = resp.Header
	logger.DebugJSON("httpclient", "resp_header", wrapper.Header)
	logger.DebugString("httpclient", "resp_body", string(wrapper.Body))
	logger.DebugString("httpclient", "resp_status", fmt.Sprint(wrapper.StatusCode))
	return wrapper
}

func addCerts(request *RequestData) error {
	if request.Cert != "" && request.PrivateKey != "" {
		clientCrt, err := tls.LoadX509KeyPair(request.Cert, request.PrivateKey)
		if err != nil {
			return err
		}
		transport := &http.Transport{
			TLSClientConfig: &tls.Config{
				Certificates:       []tls.Certificate{clientCrt},
				InsecureSkipVerify: true,
			},
		}
		httpClient.Transport = transport
		return nil
	}
	return nil
}

func createRequestError(err error) ResponseWrapper {
	errorMessage := fmt.Sprintf("创建HTTP请求错误-%s", err.Error())
	return ResponseWrapper{0, []byte(errorMessage), make(http.Header)}
}
