package _weixin

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

/*
参考：https://work.weixin.qq.com/help?person_id=1&doc_id=13376
使用方法：参考此文件的同名test文件

NOTE:
- 消息有长度限制
- 频率不超过20条/min
*/

type MsgFormat int8

const (
	MsgFormatText = iota
	MsgFormatMarkdown
)

const TextMaxLength = 2048
const ContentMaxLength = 4096

type ContentShip struct {
	Content             string   `json:"content"`               // <= 4KB
	MentionedList       []string `json:"mentioned_list"`        // 可选参数 ["Jack Ma","@all"]
	MentionedMobileList []string `json:"mentioned_mobile_list"` // 可选参数 ["13800001111","@all"]
}

// 消息结构体
type epWxMsg struct {
	MsgType  string       `json:"msgtype"`
	Text     *ContentShip `json:"text"`
	MarkDown *ContentShip `json:"markdown"`
}

// 不需要返回结果，只能观察日志来看是否推送成功
// 异步发送
func RawSendEpWxMsg(ct *ContentShip, format MsgFormat, webhookUri string, hint string) {
	var msg epWxMsg

	switch format {
	case MsgFormatText:
		msg.Text = ct
		msg.MsgType = "text"
		_len := len(msg.Text.Content)
		if _len == 0 {
			log.Printf("RawSendEpWxMsg hint:%s ------ Empty text with MsgFormatText", hint)
		} else if _len > TextMaxLength {
			log.Printf("RawSendEpWxMsg hint:%s -- Text content is too large, length:%d TextMaxLength:%d", hint, _len, TextMaxLength)
			return
		}
	case MsgFormatMarkdown:
		msg.MarkDown = ct
		msg.MsgType = "markdown"
		_len := len(msg.MarkDown.Content)
		if _len == 0 {
			log.Printf("RawSendEpWxMsg hint:%s ------ Empty content with MsgFormatMarkdown", hint)
		} else if _len > ContentMaxLength {
			log.Printf("RawSendEpWxMsg hint:%s -- Text content is too large, length:%d ContentMaxLength:%d", hint, _len, ContentMaxLength)
			return
		}
	}

	b, _ := json.Marshal(msg)

	reader := bytes.NewReader(b)
	req, _ := http.NewRequest("POST", webhookUri, reader)

	ctx, cancel := context.WithTimeout(req.Context(), time.Second*2)
	req = req.WithContext(ctx)

	go func() {
		defer cancel()
		rsp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Printf("RawSendEpWxMsg hint:%s ------ response err:%v", hint, err)
			return
		}
		b, _ = ioutil.ReadAll(rsp.Body)
		log.Printf("RawSendEpWxMsg hint:%s ------ response: %s", hint, b) // {"errcode":0,"errmsg":"ok"} 不要去反序列，不知道会不会变
	}()
}

const (
	hook = "https://weixin.xxx"
)

// NOTE: 不允许外面直接传hook地址
func SendEpWxMsgToXXX(ct *ContentShip, format MsgFormat) {
	RawSendEpWxMsg(ct, format, hook, "SendEpWxMsgToXXX")
}

/*
推送到其他群组的方法：参考SendEpWxMsgToXXX 再写一个方法，不要允许外面直传uri
*/
