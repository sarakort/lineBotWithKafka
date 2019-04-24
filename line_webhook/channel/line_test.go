package channel_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/line/line-bot-sdk-go/linebot"
)

func TestLineEventUnmarshal(t *testing.T) {
	var jsonBlob = []byte(`{"ReplyToken":"ef47d260678d4df3bf32405139bc05f1","Type":"message","Timestamp":1556093478267,"Source":{"type":"user","userId":"U35329f7ee900188c4f21d8e5a30c3c45"},"Message":{"type":"text","text":"dd"},"Joined":null,"Left":null,"Postback":null,"Beacon":null,"AccountLink":null,"Things":null,"Members":null}`)
	var  event linebot.Event

	err := event.UnmarshalJSON(jsonBlob)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%+v", event)
	// Output:
	// [{Name:Platypus Order:Monotremata} {Name:Quoll Order:Dasyuromorphia}]
}


func TestLineJsonUnmarshal(t *testing.T){
	var jsonBlob = []byte(`{"replyToken":"ef47d260678d4df3bf32405139bc05f1","type":"message","timestamp":1556090868157,"source":{"type":"user","userId":"U35329f7ee900188c4f21d8e5a30c3c45"},"message":{"id":"9747880046821","type":"text","text":"dd"}}`)
	var  event linebot.Event
	err := json.Unmarshal(jsonBlob, &event)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%+v", event)
}