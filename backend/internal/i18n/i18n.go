package i18n

import (
	"github.com/gin-gonic/gin"
)

type MessageMap struct {
	En string
	Zh string
}

var messages = map[string]MessageMap{
	"success": {
		En: "success",
		Zh: "成功",
	},
	"internal server error": {
		En: "internal server error",
		Zh: "内部服务器错误",
	},
	"invalid parameter": {
		En: "invalid parameter",
		Zh: "参数无效",
	},
	"invalid token": {
		En: "invalid token",
		Zh: "令牌无效",
	},
	"token expired": {
		En: "token expired",
		Zh: "令牌已过期",
	},
	"user not found": {
		En: "user not found",
		Zh: "用户不存在",
	},
	"user already exists": {
		En: "user already exists",
		Zh: "用户已存在",
	},
	"invalid password": {
		En: "invalid password",
		Zh: "密码无效",
	},
	"invitation invalid": {
		En: "invitation invalid",
		Zh: "邀请码无效",
	},
	"write db error": {
		En: "write db error",
		Zh: "数据库写入错误",
	},
	"read db error": {
		En: "read db error",
		Zh: "数据库读取错误",
	},
}

func Init() {
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
