package chatgptsdk

const (
	ROLE_USER      = "user"
	ROLE_SYSTEM    = "system"
	ROLE_ASSISTANT = "assistant"
)

// 消息体内容
type GPTMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// 问答响应的回答选项
type GPTMessageChoice struct {
	Index        int        `json:"index"`
	Message      GPTMessage `json:"message"`
	FinishReason string     `json:"finish_reason"`
}

// 问答响应
type GPTQuestionResp struct {
	ID      string                       `json:"id"`
	Object  string                       `json:"object"`
	Created int64                        `json:"created"`
	Model   string                       `json:"model"`
	Choices []struct{ GPTMessageChoice } `json:"choices"`
}

// 问答请求
type GPTQuestionRequest struct {
	Model    string       `json:"model"`
	Messages []GPTMessage `json:"messages"`
}
