package husky

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"time"
)

type HttpConfig struct {
	Url              string
	Method           string
	Header           http.Header
	Params           url.Values
	Timeout          time.Duration
	Data             any
	ShowLog          bool
	ResponseCallback func(response *http.Response, body []byte)
}

type HttpOptions func(config *HttpConfig)

func HttpWithUrl(url string) HttpOptions {
	return func(config *HttpConfig) {
		config.Url = url
	}
}

func HttpWithGetMethod() HttpOptions {
	return func(config *HttpConfig) {
		config.Method = http.MethodGet
	}
}

func HttpWithPostMethod() HttpOptions {
	return func(config *HttpConfig) {
		config.Method = http.MethodPost
	}
}

func HttpWithParam(key, val string) HttpOptions {
	return func(config *HttpConfig) {
		config.Params.Add(key, val)
	}
}

func HttpWithParams(params url.Values) HttpOptions {
	return func(config *HttpConfig) {
		for key := range params {
			config.Params.Add(key, params.Get(key))
		}
	}
}

func HttpWithHeader(key, val string) HttpOptions {
	return func(config *HttpConfig) {
		config.Header.Add(key, val)
	}
}

func WithJsonContentType() HttpOptions {
	return func(config *HttpConfig) {
		config.Header.Add("Content-Type", "application/json")
	}
}

func HttpWithTimeout(t time.Duration) HttpOptions {
	return func(config *HttpConfig) {
		config.Timeout = t
	}
}

func HttpWithData(data any) HttpOptions {
	return func(config *HttpConfig) {
		config.Data = data
	}
}

func HttpWithShowLog() HttpOptions {
	return func(config *HttpConfig) {
		config.ShowLog = true
	}
}

func HttpWithTraceId(traceId string) HttpOptions {
	return func(config *HttpConfig) {
		config.Header.Add("X-Trace-ID", traceId)
	}
}

func HttpWithResponse(callback func(response *http.Response, body []byte)) HttpOptions {
	return func(config *HttpConfig) {
		config.ResponseCallback = callback
	}
}

func HttpCall[R any](opts ...HttpOptions) (*R, error) {
	config := HttpConfig{
		Url:              "",
		Method:           http.MethodGet,
		Header:           map[string][]string{},
		Params:           map[string][]string{},
		Timeout:          time.Second * 60,
		Data:             nil,
		ShowLog:          false,
		ResponseCallback: nil,
	}
	for _, opt := range opts {
		opt(&config)
	}
	var requestUrl string
	if args := config.Params.Encode(); len(args) != 0 {
		requestUrl = config.Url + "?" + args
	} else {
		requestUrl = config.Url
	}
	var requestBody io.Reader
	if config.Data != nil {
		jsonData, err := json.Marshal(config.Data)
		if err != nil {
			return nil, err
		}
		requestBody = bytes.NewReader(jsonData)
	}
	client := http.Client{Timeout: config.Timeout}
	request, err := http.NewRequest(config.Method, requestUrl, requestBody)
	if err != nil {
		return nil, err
	}
	request.Header = config.Header
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	if config.ResponseCallback != nil {
		config.ResponseCallback(response, body)
	}
	if config.ShowLog {
		LogInfo("--------------------[HEADER]--------------------")
		for key := range response.Header {
			LogInfo(key, ":", response.Header.Get(key))
		}
		LogInfo("--------------------[ BODY ]--------------------")
		LogInfo(string(body))
	}

	var r R
	if err := json.Unmarshal(body, &r); err != nil {
		return nil, err
	}
	return &r, nil
}
