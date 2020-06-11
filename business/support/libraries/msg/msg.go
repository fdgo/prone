package msg

import (
	"encoding/json"
)

type Msg struct {
	Id      int32  `json:"id"`
	Body    []byte `json:"body"`
	Handler string `json:"-"`
	Message string `json:"Message"` // 这个字段是预留字段,用来区分是sns发过来的消息还是sqs发过来的消息
}

func (this *Msg) Unmarshal(v interface{}) error {
	return json.Unmarshal(this.Body, v)
}
