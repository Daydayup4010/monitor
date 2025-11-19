package utils

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"time"
)

type HttpClient struct {
	client  *resty.Client
	baseURL string
}

func CreateClient(baseURL string) *HttpClient {
	var client *HttpClient
	client = &HttpClient{
		client: resty.New().
			SetBaseURL(baseURL).
			SetTimeout(120 * time.Second).
			SetRetryCount(2).
			SetDebug(false),
		baseURL: baseURL,
	}
	return client
}

// RequestOptions 请求选项
type RequestOptions struct {
	PathParams  map[string]string
	QueryParams map[string]string
	Headers     map[string]string
	Body        interface{}
	Result      interface{}
	Error       interface{}
}

// DoRequest 执行通用请求
func (h *HttpClient) DoRequest(method, path string, opts RequestOptions) (*resty.Response, error) {
	request := h.client.R()

	// 设置路径参数
	if opts.PathParams != nil {
		request.SetPathParams(opts.PathParams)
	}

	// 设置查询参数
	if opts.QueryParams != nil {
		request.SetQueryParams(opts.QueryParams)
	}

	// 设置请求头
	if opts.Headers != nil {
		request.SetHeaders(opts.Headers)
	}

	// 设置请求体
	if opts.Body != nil {
		request.SetBody(opts.Body)
	}

	// 设置结果解析
	if opts.Result != nil {
		request.SetResult(opts.Result)
	}

	// 设置错误解析
	if opts.Error != nil {
		request.SetError(opts.Error)
	}

	// 执行请求
	var resp *resty.Response
	var err error

	switch method {
	case "GET":
		resp, err = request.Get(path)
	case "POST":
		resp, err = request.Post(path)
	case "PUT":
		resp, err = request.Put(path)
	case "DELETE":
		resp, err = request.Delete(path)
	case "PATCH":
		resp, err = request.Patch(path)
	default:
		return nil, fmt.Errorf("unsupported HTTP method: %s", method)
	}

	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	if resp.IsError() {
		return resp, fmt.Errorf("HTTP error: %s", resp.Status())
	}

	return resp, nil
}
