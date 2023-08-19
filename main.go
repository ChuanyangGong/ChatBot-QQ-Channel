package main

import (
	chatgptsdk "ChatGPTSDK"
	"QQBotSDK"
	"QQBotSDK/dto"
	"QQBotSDK/event"
	"QQBotSDK/openapi"
	"QQBotSDK/token"
	"QQBotSDK/websocket"
	"context"
	"fmt"
	"os"
)

var api *openapi.OpenAPI
var ctx context.Context
var gptClient *chatgptsdk.ChatGPT

// 调用 GPT 生成回复
func aTMessageEventHandler(event *dto.WSPayload, data *dto.WSATMessageData) error {
	if gptClient == nil {
		fmt.Println("gptClient is nil")
		return fmt.Errorf("gptClient is nil")
	}
	// 去除用户 AT 信息
	fmt.Printf("[Listening Receive] message: %v\n", data)
	question := dto.ETLMessage(data.Content)

	// 调用 GPT 生成回复
	var resp string
	// 判断 question 长度是否大于 0
	if len(question) > 0 {
		var err error
		if resp, err = gptClient.SendQuestionToGPTSimple(question); err != nil {
			resp = "出错了，请联系管理员！"
		}
	} else {
		resp = "您好我是问答小能手，你可以试着问我一些问题！"
	}
	// 发送回复
	reply, err := api.PostMessage(ctx, data.ChannelID, (&dto.PostMessage{
		Content: resp,
	}).AddAtUsr(data.Author.ID))
	if err != nil {
		fmt.Printf("[Send Message] failed, err: %v\n", err)
		return err
	}
	fmt.Printf("[reply] reply: %v", reply)
	return err
}

func main() {
	// 第一步：实例化 token 并读取配置文件
	token, err := token.CreateDefaultToken().
		ReadFromConfig("config.yaml")
	if err != nil {
		fmt.Println(err)
		return
	}
	// 第二步：实例化 GPT
	gptClient = chatgptsdk.NewClient(token.OpenaiToken, token.UseClashAsProxy)

	// 第三步：实例化 OpenAPI
	api = QQBotSDK.NewOpenAPI(token)

	// 第四步：创建 websocket 客户端并启动
	ctx := context.Background()

	// 第五步：获取 websocket 接入信息
	wsAp, err := api.GetGateway(ctx)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// 第六步：注册事件处理函数
	var atMessage event.ATMessageEventHandler = aTMessageEventHandler
	intent := event.RegisteHandlers(atMessage)

	// 第七步：创建 websocket 客户端并启动
	wsClient := websocket.CreateClient(wsAp, token, &intent)
	wsClient.Start()
}
