package _weixin

import "testing"

func TestSendEpWxMsgToBackend_Text(t *testing.T) {
	// 不需要换行符 \n
	ct := ContentShip{
		Content: `【请忽略|text|@all|@18587904107 test】 实时新增用户反馈，数据如下:
- 类型: 用户反馈
- 普通用户反馈: 117例
- VIP用户反馈: 15例`,
		MentionedList:       nil,
		MentionedMobileList: []string{"@all", "18587904107"},
	}

	SendEpWxMsgToXXX(&ct, MsgFormatText)
}

func TestSendEpWxMsgToBackend_Markdown(t *testing.T) {
	// 需要双空格换行
	ct := ContentShip{
		Content: `【请忽略|markdown|@all|@18587904107 test】实时新增用户反馈<font color="warning">132例</font>，请相关同事注意。  
>类型:<font color="comment">用户反馈</font>  
>普通用户反馈:<font color="comment">117例</font>  
>VIP用户反馈:<font color="comment">15例</font>`,

		// markdown不支持一下字段
		//MentionedList:       nil,
		//MentionedMobileList: []string{},
	}

	SendEpWxMsgToXXX(&ct, MsgFormatMarkdown)
}
