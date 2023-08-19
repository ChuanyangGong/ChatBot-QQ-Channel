# ChatBot-QQ-Channel
这是一个简易的 QQ 频道问答机器人，只需要经过简单的文件配置，即可让你的 QQ频道 拥有一个强大的问答机器人。（本项目的机器人问答能力基于 ChatGPT）

# 体验说明
本项目已部署在 [C. Y. 的校园频道](https://pd.qq.com/s/an3ew7p64) 的 **校园BBS - 聊天 | 闲聊交友** 区，可在该区域体验本项目的问答功能（注，由于 ChatGPT API 需要收费，通常并不开发，如需体验，请联系作者进行体验）。

# 使用说明
## Go 环境配置
若您的服务器上还没有 Go 开发环境，需要首先配置 Go 语言环境，详情可访问 [Golang官网](https://go.dev/)
### linux 环境配置
#### 下载 Go
在命令行输入如下命令来下载 Go 语言压缩包
```
wget https://dl.google.com/go/go1.18.10.linux-amd64.tar.gz
```

解压缩到指定目录，这样 Go 语言就安装好了
```
sudo tar -xzf go1.18.10.linux-amd64.tar.gz -C /usr/local
```

#### 修改profile
点击打开 `/etc/profile`，注意保存：

在最后面加入Go的路径:

```
export GOROOT=/usr/local/go
export PATH=$PATH:$GOROOT/bin 
```

使profile文件生效
```
source /etc/profile
```

在命令行输入 `go version` 指令检验是否安装完成，如果安装成功，会打印出 Go 的版本号

## 申请 QQ 机器人账号
在官网申请机器人账号，申请完可在机器人后台 - 开发 - 开发设置中找到 `BotAppID` 和 `机器人令牌`
- [机器人申请指南](https://bot.q.qq.com/wiki/#%E7%AE%80%E4%BB%8B)
- [机器人后台官网](https://q.qq.com/#/)


## 申请 ChatGPT API Token
在 [ChatGPT 官网](https://platform.openai.com/account/api-keys)申请账号，并申请 `API Key` 
```
sk-xxxxxxxxxxxxxxxx
```

## 克隆项目
将本项目克隆至服务器
```
git clone https://github.com/ChuanyangGong/QQ-Channel-ChatBot.git
```
进入项目目录并拷贝一份配置文件
```
cd QQ-Channel-ChatBot
cp config.yaml.example config.yaml
```
修改配置文件
```
vim config.yaml
```
将配置内容写入，其中 `use_clash_as_proxy` 用于使用 clash 进行 chatgpt openapi 的代理，默认配置为了代理到 Clash 的默认端口 `7890`，不进行代理可能无法连接到 OpenAI 的服务器
```
appid: 601xxxxxx                  # BotAppID
token: Anrxxxxxxxxxxxx            # 机器人令牌
openai_token: sk-xxxxxxxxxxxxxx   # ChatGPT API Key
use_clash_as_proxy: false         # 是否使用 Clash 进行代理
```

## 运行项目
```
go run main.go
```

# 项目设计
本项目主要涉及两个大的模块，分别是 `ChatGPTSDK` 和 `QQBotSDK`。

## ChatGPTSDK
`ChatGPTSDK` 用于与 OpenAI 服务器进行连接，封装了与 OpenAI 进行交互的细节，提供了 `ChatGPTSDK.NewClient()` 方法用于创建一个连接实例，以及 `SendQuestionToGPTSimple()` 方法用于向 OpenAI 服务器发送问题，并获取回答。

## QQBotSDK
`QQBotSDK` 封装了 [QQ机器人平台的接口](https://bot.q.qq.com/wiki/develop/api/#%E6%8E%A5%E5%8F%A3%E8%AF%B4%E6%98%8E)，主要提供 `openapi` 和 `websocket` 两种实例的获取方法。

- openapi 上封装了该项目会用到的 HTTP 请求的接口，包括 [获取通用 WSS 接入点](https://bot.q.qq.com/wiki/develop/api/openapi/wss/url_get.html) 、[发送消息](https://bot.q.qq.com/wiki/develop/api/openapi/message/post_messages.html) 等接口，根据需要可继续扩展

- websocket 上封装了与 QQ机器人进行 websocket 连接的细节，会自动进行鉴权、心跳等操作，用户自定义的事件可通过 `QQBotSDK.event.RegisteHandlers()` 进行注册，目前已实现了对 `AT_MESSAGE_CREATE` 事件的处理，支持进一步扩展

# 联系方式

- email: cy_gong@foxmail.com