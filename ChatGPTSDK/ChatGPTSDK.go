package chatgptsdk

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/go-resty/resty/v2"
)

const GPT_TOKEN_TYPE = "Bearer"
const GPT_Model = "gpt-3.5-turbo"
const CLASH_PROXY = "http://127.0.0.1:7890"
const OPENAI_API_URL = "https://api.openai.com/v1/chat/completions"

type ChatGPT struct {
	AuthToken    string
	restryClient *resty.Client // resty 客户端
}

// 创建 OpenAPI 实例
func NewClient(authToken string, useProxy bool) *ChatGPT {
	api := &ChatGPT{
		AuthToken:    authToken,
		restryClient: resty.New(),
	}
	api.initRestyClient(useProxy)
	return api
}

// 初始化 resty client
func (api *ChatGPT) initRestyClient(useProxy bool) {
	if api.restryClient == nil {
		api.restryClient = resty.New()
	}
	proxyURL, err := url.Parse(CLASH_PROXY)
	if err != nil {
		panic(err)
	}
	if useProxy {
		api.restryClient.SetTransport(&http.Transport{
			Proxy: http.ProxyURL(
				proxyURL,
			),
		})
	}
	api.restryClient.SetAuthToken(api.AuthToken).
		SetAuthScheme(GPT_TOKEN_TYPE).
		OnAfterResponse(
			func(client *resty.Client, response *resty.Response) error {
				fmt.Println(respInfo(response))
				return nil
			},
		)
}

// 发送一个请求
func (api *ChatGPT) request(ctx context.Context) *resty.Request {
	return api.restryClient.R().SetContext(ctx)
}

// 简易版的发送消息请求
func (api *ChatGPT) SendQuestionToGPTSimple(question string) (string, error) {
	question = "从现在开始，你的名字叫做问答小能手，可以解答我问你的任何问题，如果我问你你是谁，你能干什么，" +
		"你就告诉我你叫问答小能手，可以解答我的问题。现在请回答我的问题：" + question
	resp, err := api.SendQuestionToGPT(question)
	if err != nil || len(resp.Choices) == 0 {
		return "", err
	}
	return resp.Choices[0].Message.Content, nil
}

// 发送消息请求
func (api *ChatGPT) SendQuestionToGPT(question string) (*GPTQuestionResp, error) {
	resp, err := api.request(context.Background()).
		SetResult(GPTQuestionResp{}).
		SetBody(GPTQuestionRequest{
			Model: GPT_Model,
			Messages: []GPTMessage{
				{
					Role:    ROLE_USER,
					Content: question,
				},
			},
		}).
		Post(OPENAI_API_URL)
	if err != nil {
		return nil, err
	}

	return resp.Result().(*GPTQuestionResp), nil
}

// 格式化请求/响应参数
func respInfo(resp *resty.Response) string {
	bodyJson, _ := json.Marshal(resp.Request.Body)
	return fmt.Sprintf(
		"[OPENAPI] method: %s, url: %s, status: %v, req: %v, resp: %v",
		resp.Request.Method,
		resp.Request.URL,
		resp.Status(),
		string(bodyJson),
		string(resp.Body()),
	)
}
