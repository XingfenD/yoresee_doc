package i18n

import (
	"github.com/XingfenD/yoresee_doc/internal/status"
	"github.com/gin-gonic/gin"
)

type MessageMap struct {
	En string
	Zh string
}

var messages = make(map[string]MessageMap)

// zhTranslations 手动维护的中文翻译，使用status包中的变量作为键
var zhTranslations = map[string]string{
	// 直接使用status包中的消息作为键
	status.StatusSuccess.(*status.Status).Message:                       "成功",
	status.StatusParamError.(*status.Status).Message:                    "参数无效",
	status.StatusTokenInvalid.(*status.Status).Message:                  "令牌无效",
	status.StatusTokenExpired.(*status.Status).Message:                  "令牌已过期",
	status.StatusUserNotFound.(*status.Status).Message:                  "用户不存在",
	status.StatusUserAlreadyExists.(*status.Status).Message:             "用户已存在",
	status.StatusInvalidPassword.(*status.Status).Message:               "密码无效",
	status.StatusInvitationInvalid.(*status.Status).Message:             "邀请码无效",
	status.StatusPermissionDenied_DocumentRead.(*status.Status).Message: "没有文档读取权限",
	status.StatusWriteDBError.(*status.Status).Message:                  "数据库写入错误",
	status.StatusReadDBError.(*status.Status).Message:                   "数据库读取错误",
}

// Init 初始化i18n消息，从status包加载所有状态消息
func Init() {
	// 从status包加载所有状态消息
	loadStatusMessages()
}

// loadStatusMessages 从zhTranslations加载所有状态消息
func loadStatusMessages() {
	// 直接从zhTranslations加载所有消息
	// zhTranslations的键已经是status消息，值是中文翻译
	for msg, zh := range zhTranslations {
		messages[msg] = MessageMap{
			En: msg, // 使用消息本身作为英文
			Zh: zh,  // 使用预定义的中文翻译
		}
	}
}

func Translate(c *gin.Context, key string) string {
	lang := c.GetHeader("Accept-Language")
	if lang == "" {
		lang = "en"
	}

	if msg, ok := messages[key]; ok {
		if lang == "zh" {
			return msg.Zh
		}
		return msg.En
	}

	return key
}

// GetMessageMap 获取消息映射，用于测试或调试
func GetMessageMap() map[string]MessageMap {
	return messages
}
