package channel_test

import (
	"encoding/json"
	"testing"
	"fmt"

	"bot/connect/channel"
)

func TestBotApiUnmarshal(t *testing.T){
	var jsonBlob = []byte(`{"topic": "\u0e02\u0e2d\u0e23\u0e32\u0e04\u0e32\u0e2a\u0e48\u0e07", "confidence": "60%"}`)
	var reply channel.BotapiReply

	if err := json.Unmarshal(jsonBlob, &reply); err != nil {
		t.Error(err)
	}

	fmt.Printf("%+v", reply)
}
