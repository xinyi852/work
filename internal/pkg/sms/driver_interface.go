package sms

type Driver interface {
	// 参数短信
	Send(phone string, message Message, config map[string]string) bool
}
