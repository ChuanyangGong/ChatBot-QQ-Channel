package main

import (
	"QQBotSDK"
	"QQBotSDK/dto"
	"QQBotSDK/event"
	"QQBotSDK/openapi"
	"QQBotSDK/token"
	"QQBotSDK/websocket"
	"context"
	"fmt"
	"os"
	"strings"
)

var api *openapi.OpenAPI
var ctx context.Context

func aTMessageEventHandler(event *dto.WSPayload, data *dto.WSATMessageData) error {
	fmt.Printf("[Listening Receive] message: %v\n", data)
	if strings.HasSuffix(data.Content, "> hello") {
		ctx = context.Background()
		reply, err := api.PostMessage(ctx, data.ChannelID, (&dto.PostMessage{
			Content: "你好",
		}).AddAtUsr(data.Author.ID))
		if err != nil {
			fmt.Printf("[Send Message] failed, err: %v\n", err)
			return err
		}
		fmt.Printf("[reply] reply: %v", reply)
	}
	return nil
}

func main() {
	// 第一步：实例化 token 并读取配置文件
	token, err := token.CreateDefaultToken().
		ReadFromConfig("config.yaml")
	if err != nil {
		fmt.Println(err)
		return
	}

	// 第二步：实例化 OpenAPI
	api = QQBotSDK.NewOpenAPI(token)
	// 第三步：初始化 context
	ctx := context.Background()
	// 第四步：获取 websocket 接入信息
	wsAp, err := api.GetGateway(ctx)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var atMessage event.ATMessageEventHandler = aTMessageEventHandler
	intent := event.RegisteHandlers(atMessage)

	wsClient := websocket.CreateClient(wsAp, token, &intent)
	wsClient.Start()
}
