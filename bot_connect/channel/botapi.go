package channel


type BotapiReply struct {
	ReplyToken string `json:"replytoken"`
	Topic string `json:"topic"`
	Confidence string `json:"confidence"`
}