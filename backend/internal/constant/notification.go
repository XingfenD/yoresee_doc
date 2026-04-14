package constant

// Notification type constants.
const (
	NotificationType_Comment = "comment"
	NotificationType_Reply   = "reply"
	NotificationType_Mention = "mention"
	NotificationType_System  = "system"
)

// notificationTitles holds default (zh-CN) titles keyed by notification type.
// Replace with a proper i18n lookup once multi-language support is added.
var notificationTitles = map[string]string{
	NotificationType_Comment: "收到新评论",
	NotificationType_Reply:   "回复了你的评论",
	NotificationType_Mention: "在评论中提及了你",
	NotificationType_System:  "系统通知",
}

// GetNotificationTitle returns the default title for the given notification type.
// Falls back to the type string itself if the type is unrecognised.
func GetNotificationTitle(notifType string) string {
	if title, ok := notificationTitles[notifType]; ok {
		return title
	}
	return notifType
}
