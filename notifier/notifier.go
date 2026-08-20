package notifier

import (
	"fmt"
)

type Notifier interface {
	Send(message string)
}

type EmailNotifier struct {
	Email string
}

func (e EmailNotifier) Send(message string) {
	fmt.Printf("📧 [EMAIL отправлен на %s]: %s\n", e.Email, message)
}

type TelegramNotifier struct {
	Username string
}

func (t TelegramNotifier) Send(message string) {
	fmt.Printf("✈️ [TELEGRAM отправлен пользователю @%s]: %s\n", t.Username, message)
}
